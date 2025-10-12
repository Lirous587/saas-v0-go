package service

import (
	"fmt"
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

func (s *service) Create(comment *domain.Comment, belongKey string) (*domain.Comment, error) {
	// 1.plate 是否存在
	plate, err := s.repo.GetPlateBelongByKey(comment.TenantID, belongKey)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	comment.PlateID = plate.ID

	// 2.验证root_id和parent_id合理性
	// 当前板块下是否存在root_id和parent_id
	if comment.RootID != 0 {
		exist, err := s.repo.IsCommentInPlate(comment.TenantID, comment.PlateID, comment.RootID)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if !exist {
			return nil, codes.ErrCommentNotFoundInNowPlate
		}
	}
	if comment.ParentID != 0 {
		exist, err := s.repo.IsCommentInPlate(comment.TenantID, comment.PlateID, comment.ParentID)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if !exist {
			return nil, codes.ErrCommentNotFoundInNowPlate
		}
	}

	// 3.创建评论
	comment, err = s.repo.Create(comment)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// 4.异步发送邮件通知
	// 判断评论者角色 domain_admin viewer
	// domain_admin 评论 如果root parent存在且不为自己 才需发送任何邮件
	// viewer 评论 如果root parent为自己 也无需发送邮件
	go func() {
		// 查询租户管理员id
		adminId, err := s.repo.GetDomainAdminByTenant(comment.TenantID)
		if err != nil {
			zap.L().Error("获取租户管理用户失败",
				zap.Int64("tenant_id", int64(comment.TenantID)),
				zap.Int64("comment_id", comment.ID),
				zap.Error(err))
			return
		}

		var toUserIds []int64

		// 获取toUserIds
		if comment.IsReply() {
			// 回复
			// 首先根据root_id parent_id拿到对应的userIds
			// domain_admin回复 过滤自己之后 给filteredUserIds发邮件
			// viewer回复 过滤自己之后 给filteredUserIds以及domain_admin邮件
			// 都需要去重处理 避免多次收到邮件
			userIds, err := s.repo.GetUserIdsByRootORParent(comment.TenantID, comment.PlateID, comment.RootID, comment.ParentID)
			if err != nil {
				zap.L().Error("获取根/父评论用户失败",
					zap.Int64("tenant_id", int64(comment.TenantID)),
					zap.Int64("comment_id", comment.ID),
					zap.Int64("root_id", comment.RootID),
					zap.Int64("parent_id", comment.ParentID),
					zap.Error(err))
				return
			}

			fmt.Println("userIds", userIds)

			// 从userIds中排除自己
			filteredIds := comment.FilterSelf(userIds)

			fmt.Println("filteredIds", filteredIds)

			if comment.IsCommentByAdmin(adminId) {
				toUserIds = utils.UniqueInt64s(filteredIds)
			} else {
				// 要加上admin_id
				toUserIds = utils.UniqueInt64s(append(filteredIds, adminId))
			}

		} else {
			// 创建根评论
			// domain_admin创建无需发送邮件
			// viewer创建需要给domain_admin发送邮件
			if comment.IsCommentByAdmin(adminId) {
				zap.L().Info("创建根评论，评论用户为租户管理员，无需发送邮件")
			} else {
				toUserIds = []int64{adminId}
			}
		}

		if len(toUserIds) == 0 {
			zap.L().Info("无用户需要发送邮件",
				zap.Int64("tenant_id", int64(comment.TenantID)),
				zap.Int64("comment_id", comment.ID))
			return
		}

		// 获取当前评论用户
		commentUser, err := s.repo.GetUserInfoByID(comment.UserID)
		if err != nil {
			zap.L().Error("获取当前评论用户失败",
				zap.Int64("comment_id", comment.ID),
				zap.Int64("user_id", comment.UserID),
				zap.Error(err))
			return
		}

		// 查询所要发送邮件的用户
		toUsers, err := s.repo.GetUserInfosByIds(toUserIds)
		if err != nil {
			zap.L().Error("获取用户信息失败",
				zap.Int64("tenant_id", int64(comment.TenantID)),
				zap.Int64("comment_id", comment.ID),
				zap.Int64s("user_ids", toUserIds),
				zap.Error(err))
			return
		}

		// 获取 板块的related_url
		relatedURL, err := s.repo.GetPlateRelatedURlByID(comment.TenantID, plate.ID)
		if err != nil {
			zap.L().Error("获取板块RelatedURl失败",
				zap.Int64("tenant_id", int64(comment.TenantID)),
				zap.Int64("comment_id", comment.ID),
				zap.Int64("plate_key", comment.PlateID),
				zap.Error(err))
			return
		}

		// 限制 goroutine 数量
		sem := make(chan struct{}, 10) // 最多 10 个并发
		for _, toUser := range toUsers {
			go func(u *domain.UserInfo) {
				sem <- struct{}{}        // 获取信号
				defer func() { <-sem }() // 释放
				if err := s.sentCommentEmail(commentUser, u.GetEmail(), relatedURL, comment.Content); err != nil {
					zap.L().Error("发送邮件失败",
						zap.Int64("tenant_id", int64(comment.TenantID)),
						zap.Int64("comment_id", comment.ID),
						zap.Int64("to_user_id", u.ID),
						zap.String("to_email", u.GetEmail()),
						zap.Error(err))
					return
				}
			}(toUser)
		}
	}()

	return comment, nil
}

func (s *service) Delete(tenantID domain.TenantID, userID int64, id int64) error {
	// 查询当前评论用户
	uid, err := s.repo.GetCommentUser(tenantID, id)
	if err != nil {
		return errors.WithStack(err)
	}

	// 如果请求用户和评论用户不一致
	if uid != userID {
		// 去获取当前租户的uid
		adminID, err := s.repo.GetDomainAdminByTenant(tenantID)
		if err != nil {
			return errors.WithStack(err)
		}

		if userID != adminID {
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

func (s *service) DeletePlate(tenantID domain.TenantID, id int64) error {
	return s.repo.DeletePlate(tenantID, id)
}

func (s *service) ListPlate(query *domain.PlateQuery) (*domain.PlateList, error) {
	return s.repo.ListPlate(query)
}

func (s *service) SetTenantConfig(config *domain.TenantConfig) error {
	// 生成client_token
	clientToken, err := utils.GenRandomHexToken()
	if err != nil {
		return err
	}

	config.ClientToken = clientToken

	return s.repo.SetTenantConfig(config)
}

func (s *service) GetTenantConfig(tenantID domain.TenantID) (*domain.TenantConfig, error) {
	return s.repo.GetTenantConfig(tenantID)
}

func (s *service) SetPlateConfig(config *domain.PlateConfig) error {
	plate, err := s.repo.GetPlateBelongByKey(config.TenantID, config.Plate.BelongKey)
	if err != nil {
		return errors.WithStack(err)
	}

	config.Plate.ID = plate.ID
	if err := s.repo.SetPlateConfig(config); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *service) GetPlateConfig(tenantID domain.TenantID, belongKey string) (*domain.PlateConfig, error) {
	plate, err := s.repo.GetPlateBelongByKey(tenantID, belongKey)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	config, err := s.repo.GetPlateConfig(tenantID, plate.ID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return config, nil
}
