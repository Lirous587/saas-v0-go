package service

import (
	"saas/internal/comment/domain"
	"saas/internal/common/email"
	"saas/internal/common/reskit/codes"
	"saas/internal/common/utils"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
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

func (s *service) Audit(tenantID domain.TenantID, commentID domain.CommentID, status domain.CommentStatus) error {
	comment, err := s.repo.GetByID(tenantID, commentID)
	if err != nil {
		return errors.WithStack(err)
	}

	if !comment.CanAudit() {
		return codes.ErrCommentIllegalAudit
	}

	if status.IsApproved() {
		comment.SetApproved()
	}

	// 同意
	if comment.IsApproved() {
		if err := s.repo.Approve(tenantID, commentID); err != nil {
			return errors.WithMessage(err, "同意评论时候更新status失败")
		}
	} else {
		if err := s.repo.Delete(tenantID, commentID); err != nil {
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
				uids, err := s.repo.GetUserIDsByRootORParent(tenantID, comment.PlateID, comment.RootID, comment.ParentID)
				if err != nil {
					zap.L().Error("根据当前评论的root和parent去查询uids失败", zap.Error(err))
					return
				}

				// 2.查询租户
				admin, err := s.repo.GetTenantCreator(tenantID)
				if err != nil {
					zap.L().Error("获取租户管理员失败", zap.Error(err))
					return
				}

				// 通知其回复人员 从uids中除去自己和租户管理员(因为此时租户审核了就无需通知)
				filterSelfIDs := comment.FilterSelf(uids)
				filterUIDs := make([]domain.UserID, 0, 3)
				for _, id := range filterSelfIDs {
					if id == admin.ID {
						continue
					}
					filterUIDs = append(filterUIDs, id)
				}

				filterUIDsStr := domain.UserIDs(filterUIDs).ToStringSlice()

				toUIDs := utils.UniqueStrings(filterUIDsStr)

				toUserIDs := domain.NewUserIDsFromStrings(toUIDs)

				// 获取待通知的用户信息
				toUsers, err := s.repo.GetUserInfosByIDs(toUserIDs)
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

func (s *service) Delete(tenantID domain.TenantID, userID domain.UserID, commentID domain.CommentID) error {
	// 查询当前评论用户
	uid, err := s.repo.GetCommentUser(tenantID, commentID)
	if err != nil {
		return errors.WithStack(err)
	}

	// 如果请求用户和评论用户不一致
	if uid != userID {
		// 获取当前租户管理员
		admin, err := s.repo.GetTenantCreator(tenantID)
		if err != nil {
			return errors.WithStack(err)
		}

		if userID != admin.ID {
			return codes.ErrCommentNoPermissionToDelete
		}
	}

	return s.repo.Delete(tenantID, commentID)
}

type commentLikeHelper interface {
	CommentID() domain.CommentID
	Like()
}

// 泛型函数：处理点赞记录
func applyLikeStatus[T commentLikeHelper](repo domain.CommentRepository, comments []T, tenantID domain.TenantID, userID domain.UserID) {
	if userID.IsZero() || comments == nil {
		zap.L().Debug("用户未登录或无评论数据 applyLikeStatus无效处理后续逻辑")
		return
	}

	commentIDs := make([]domain.CommentID, 0, len(comments))

	for i := range comments {
		commentIDs = append(commentIDs, comments[i].CommentID())
	}

	ids, err := repo.GetLikeRecords(tenantID, commentIDs, userID)
	if err != nil {
		zap.L().Error("获取点赞状态失败", zap.Error(err))
	}

	if len(ids) == 0 {
		return
	}

	cMap := domain.CommentIDs(ids).ToMap()

	for i := range comments {
		_, exist := cMap[comments[i].CommentID()]
		if exist {
			comments[i].Like()
		}
	}
}

func (s *service) ListRoots(belongKey string, userID domain.UserID, query *domain.CommentRootsQuery) ([]*domain.CommentRoot, error) {
	plateID, err := s.getPlateID(query.TenantID, belongKey)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	query.PlateID = plateID

	// 获取评论基础信息和相关用户信息
	roots, err := s.repo.ListRoots(query)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// 为评论设置点赞状态
	applyLikeStatus(s.repo, roots, query.TenantID, userID)

	return roots, nil
}

func (s *service) ListReplies(belongKey string, userID domain.UserID, query *domain.CommentRepliesQuery) ([]*domain.CommentReply, error) {
	plateID, err := s.getPlateID(query.TenantID, belongKey)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	query.PlateID = plateID

	// 获取到replies
	replies, err := s.repo.ListReplies(query)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	applyLikeStatus(s.repo, replies, query.TenantID, userID)

	return replies, nil
}

func (s *service) ListNoAudits(belongKey string, query *domain.CommentNoAuditQuery) ([]*domain.CommentNoAudit, error) {
	plateID, err := s.getPlateID(query.TenantID, belongKey)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	query.PlateID = plateID

	return s.repo.ListNoAudits(query)
}

func (s *service) ToggleLike(tenantID domain.TenantID, userID domain.UserID, commentID domain.CommentID) error {
	// 去查询当前status
	likeStatus, err := s.repo.GetLikeStatus(tenantID, commentID, userID)
	if err != nil {
		return errors.WithStack(err)
	}

	// toogle状态
	likeStatus.Toogle()

	var eg errgroup.Group

	// 并发执行数据库操作和点赞计数更新
	if likeStatus.IsLike() {
		eg.Go(func() error {
			return s.repo.AddLike(tenantID, commentID, userID)
		})
	} else {
		eg.Go(func() error {
			return s.repo.RemoveLike(tenantID, commentID, userID)
		})
	}

	eg.Go(func() error {
		return s.repo.UpdateLikeCount(tenantID, commentID, likeStatus.IsLike())
	})

	return eg.Wait()
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

	// 删除缓存
	if err := s.cache.DeletePlateID(plate.TenantID, plate.BelongKey); err != nil {
		zap.L().Error(
			"删除板块ID缓存失败",
			zap.Error(err),
			zap.String("tenant_id", plate.TenantID.String()),
			zap.String("belong_key", plate.BelongKey),
		)
	}

	return nil
}

func (s *service) DeletePlate(tenantID domain.TenantID, plateID domain.PlateID) error {
	return s.repo.DeletePlate(tenantID, plateID)
}

func (s *service) ListPlate(query *domain.PlateQuery) (*domain.PlateList, error) {
	return s.repo.ListPlate(query)
}

func (s *service) getPlateID(tenantID domain.TenantID, belongKey string) (domain.PlateID, error) {
	// 尝试从缓存获取
	plateID, cacheErr := s.cache.GetPlateID(tenantID, belongKey)
	if cacheErr == nil {
		return plateID, nil
	}

	// 缓存未命中或出错，从数据库获取
	belong, err := s.repo.GetPlateBelongByKey(tenantID, belongKey)
	if err != nil {
		return "", errors.WithStack(err)
	}

	// 如果是缓存缺失，异步写入缓存
	if errors.Is(cacheErr, codes.ErrCommentPlateIDCacheMissing) {
		go func() {
			if setErr := s.cache.SetPlateID(tenantID, belong.BelongKey, belong.ID); setErr != nil {
				zap.L().Error(
					"设置板块ID缓存失败",
					zap.Error(setErr),
					zap.String("tenant_id", tenantID.String()),
					zap.String("belong_key", belong.BelongKey),
					zap.String("plate_id", belong.ID.String()),
				)
			}
		}()
	}

	return belong.ID, nil
}

func (s *service) SetTenantConfig(config *domain.TenantConfig) error {
	// 删除缓存
	if err := s.cache.DeleteTenantConfig(config.TenantID); err != nil {
		zap.L().Error(
			"删除租户级别评论配置缓存失败",
			zap.Error(err),
			zap.String("tenant_id", config.TenantID.String()),
		)
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
					zap.String("tenant_id", tenantID.String()),
				)
			}
		}()
	}

	return config, nil
}

func (s *service) SetPlateConfig(config *domain.PlateConfig) error {
	plateID, err := s.getPlateID(config.TenantID, config.Plate.BelongKey)
	if err != nil {
		return errors.WithStack(err)
	}
	config.Plate.ID = plateID

	// 删除缓存
	if err := s.cache.DeletePlateConfig(config.TenantID, config.Plate.ID); err != nil {
		zap.L().Error(
			"删除板块级别评论配置缓存失败",
			zap.Error(err),
			zap.String("tenant_id", config.TenantID.String()),
			zap.String("plate_id", config.Plate.ID.String()),
		)
	}

	if err := s.repo.SetPlateConfig(config); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *service) GetPlateConfig(tenantID domain.TenantID, plateID domain.PlateID) (*domain.PlateConfig, error) {
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
					zap.String("tenant_id", tenantID.String()),
					zap.String("plate_id", plateID.String()),
				)
			}
		}()
	}

	return config, nil
}

func (s *service) CheckPlateBelongKey(tenantID domain.TenantID, belongKey string) (bool, error) {
	return s.repo.ExistPlateBykey(tenantID, belongKey)
}
