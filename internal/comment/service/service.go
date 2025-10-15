package service

import (
	"saas/internal/comment/domain"
	"saas/internal/common/email"
	"saas/internal/common/reskit/codes"
	"saas/internal/common/utils"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type service struct {
	repo   domain.CommentRepository
	cache  domain.CommentCache
	mailer email.Mailer
}

func NewCommentService(repo domain.CommentRepository, cache domain.CommentCache, mailer email.Mailer) domain.CommentService {
	return &service{
		repo:   repo,
		cache:  cache,
		mailer: mailer,
	}
}

func (s *service) Audit(tenantID domain.TenantID, id int64, status domain.CommentStatus) error {
	comment, err := s.repo.GetByID(tenantID, id)
	if err != nil {
		return errors.WithStack(err)
	}

	if !comment.CanAudit() {
		return codes.ErrCommentIllegalAudit
	}

	if status == domain.CommentStatusApproved {
		comment.SetApproved()
	}

	// 同意
	if comment.IsApproved() {
		if err := s.repo.Approve(tenantID, id); err != nil {
			return errors.WithMessage(err, "同意评论时候更新status失败")
		}
	} else {
		if err := s.repo.Delete(tenantID, id); err != nil {
			return errors.WithMessage(err, "拒绝评论时候删除评论记录失败")
		}
	}

	go func() {
		// 获取评论来源
		commentSource, err := s.getCommentSource(comment)
		if err != nil {
			zap.L().Error("获取评论来源失败", zap.Error(err))
			return
		}

		// 后续异步处理
		// 此处与上方拆开 避免逻辑混乱
		if comment.IsApproved() {
			go func() {
				// - 通知评论者
				// 通知评论用户(必定不为租户管理员 放心处理)
				if err := s.sentAuditPassEmail(commentSource.commentUser.GetEmail(), commentSource.relatedURL, comment.Content); err != nil {
					zap.L().Error("发送邮件AuditPass失败", zap.Error(err))
					return
				}

				// - 通知回复者
				// 1.根据当前评论的root和parent去查询uids
				uids, err := s.repo.GetUserIdsByRootORParent(tenantID, comment.PlateID, comment.RootID, comment.ParentID)
				if err != nil {
					zap.L().Error("根据当前评论的root和parent去查询uids失败", zap.Error(err))
					return
				}

				// 2.查询租户
				admin, err := s.repo.GetDomainAdminByTenant(tenantID)
				if err != nil {
					zap.L().Error("获取租户管理员失败", zap.Error(err))
					return
				}

				// 通知其回复人员 从uids中除去自己和租户管理员(因为此时租户审核了就无需通知)
				filterSelfIds := comment.FilterSelf(uids)
				filterIds := make([]int64, 0, 3)
				for _, id := range filterSelfIds {
					if id == admin.ID {
						continue
					}
					filterIds = append(filterIds, id)
				}

				toUids := utils.UniqueInt64s(filterIds)
				// 获取待通知的用户信息
				toUsers, err := s.repo.GetUserInfosByIds(toUids)
				if err != nil {
					zap.L().Error("获取待通知的用户信息失败", zap.Error(err))
					return
				}
				// 整合数据 发送邮件
				for _, toUser := range toUsers {
					go func(u *domain.UserInfo) {
						if err := s.sentCommentEmail(commentSource.commentUser, u.GetEmail(), commentSource.relatedURL, comment.Content); err != nil {
							zap.L().Error("发送邮件commentEmail失败", zap.Error(err))
							return
						}
					}(toUser)
				}

			}()
		} else {
			go func() {
				if err := s.sentAuditRejectEmail(commentSource.commentUser.GetEmail(), commentSource.relatedURL, comment.Content); err != nil {
					zap.L().Error("发送邮件auditRejectEmail失败", zap.Error(err))
				}
			}()
		}
	}()

	return nil
}

func (s *service) Delete(tenantID domain.TenantID, userID int64, id int64) error {
	// 查询当前评论用户
	uid, err := s.repo.GetCommentUser(tenantID, id)
	if err != nil {
		return errors.WithStack(err)
	}

	// 如果请求用户和评论用户不一致
	if uid != userID {
		// 获取当前租户管理员
		admin, err := s.repo.GetDomainAdminByTenant(tenantID)
		if err != nil {
			return errors.WithStack(err)
		}

		if userID != admin.ID {
			return codes.ErrCommentNoPermissionToDelete
		}
	}

	return s.repo.Delete(tenantID, id)
}

func (s *service) List(query *domain.CommentQuery) (*domain.CommentList, error) {
	return s.repo.List(query)
}

func (s *service) CreatePlate(plate *domain.Plate) error {
	exist, err := s.repo.ExistPlateBykey(plate.TenantID, plate.BelongKey)
	if err != nil {
		return errors.WithStack(err)
	}

	if exist {
		return codes.ErrCommentPlateExist
	}

	if err := s.repo.CreatePlate(plate); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s *service) UpdatePlate(plate *domain.Plate) error {

	if err := s.repo.UpdatePlate(plate); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s *service) DeletePlate(tenantID domain.TenantID, id int64) error {
	return s.repo.DeletePlate(tenantID, id)
}

func (s *service) ListPlate(query *domain.PlateQuery) (*domain.PlateList, error) {
	return s.repo.ListPlate(query)
}

func (s *service) SetTenantConfig(config *domain.TenantConfig) error {
	// 判断是否已有配置
	exist, err := s.repo.ExistTenantConfigByID(config.TenantID)
	if err != nil {
		return err
	}

	// 删除缓存
	if err := s.cache.DeleteTenantConfig(config.TenantID); err != nil {
		zap.L().Error(
			"删除租户级别评论配置缓存失败",
			zap.Error(err),
			zap.Int64("tenant_id", int64(config.TenantID)),
		)
	}

	// 没配置过就生成client_token
	if !exist {
		// 生成client_token
		clientToken, err := utils.GenRandomHexToken()
		if err != nil {
			return err
		}
		config.SetClientToken(clientToken)
	}

	return s.repo.SetTenantConfig(config)
}

func (s *service) GetTenantConfig(tenantID domain.TenantID) (*domain.TenantConfig, error) {
	// 尝试从缓存获取
	config, cacheErr := s.cache.GetTenantConfig(tenantID)
	if cacheErr == nil {
		return config, nil
	}

	// 缓存未命中或出错，从数据库获取
	config, err := s.repo.GetTenantConfig(tenantID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// 如果是缓存缺失，异步写入缓存
	if errors.Is(cacheErr, codes.ErrCommentTenantConfigCacheMissing) {
		go func() {
			if setErr := s.cache.SetTenantConfig(config); setErr != nil {
				zap.L().Error(
					"设置租户级别评论配置缓存失败",
					zap.Error(setErr),
					zap.Int64("tenant_id", int64(tenantID)),
				)
			}
		}()
	}

	return config, nil
}

func (s *service) SetPlateConfig(config *domain.PlateConfig) error {
	plate, err := s.repo.GetPlateBelongByKey(config.TenantID, config.Plate.BelongKey)
	if err != nil {
		return errors.WithStack(err)
	}
	config.Plate.ID = plate.ID

	// 删除缓存
	if err := s.cache.DeletePlateConfig(config.TenantID, config.Plate.ID); err != nil {
		zap.L().Error(
			"删除板块级别评论配置缓存失败",
			zap.Error(err),
			zap.Int64("tenant_id", int64(config.TenantID)),
			zap.Int64("plate_id", config.Plate.ID),
		)
	}

	if err := s.repo.SetPlateConfig(config); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *service) GetPlateConfig(tenantID domain.TenantID, plateID int64) (*domain.PlateConfig, error) {
	// 尝试从缓存获取
	config, cacheErr := s.cache.GetPlateConfig(tenantID, plateID)
	if cacheErr == nil {
		return config, nil
	}

	// 缓存未命中或出错，从数据库获取
	config, err := s.repo.GetPlateConfig(tenantID, plateID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// 如果是缓存缺失，异步写入缓存
	if errors.Is(cacheErr, codes.ErrCommentPlateConfigCacheMissing) {
		go func() {
			if setErr := s.cache.SetPlateConfig(config); setErr != nil {
				zap.L().Error(
					"设置板块级别评论配置缓存失败",
					zap.Error(setErr),
					zap.Int64("tenant_id", int64(tenantID)),
					zap.Int64("plate_id", plateID),
				)
			}
		}()
	}

	return config, nil
}
