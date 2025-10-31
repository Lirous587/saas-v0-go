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

func (s *service) ListRoots(belongKey string, userID domain.UserID, query *domain.CommentRootsQuery) ([]*domain.CommentRoot, error) {
	plateID, err := s.getPlateID(query.TenantID, belongKey)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	query.PlateID = plateID

	// 获取到roots
	roots, err := s.repo.ListRoots(query)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// 登录用户则处理点赞
	if !userID.IsZero() {
		rootIDs := make([]domain.CommentID, 0, len(roots))
		for i := range roots {
			rootIDs = append(rootIDs, roots[i].CommentWithUser.ID)
		}

		// 整合点赞记录
		likeMap, err := s.cache.GetLikeMap(query.TenantID, userID, rootIDs)
		if err != nil {
			zap.L().Error("获取点赞状态失败", zap.Error(err))
		}

		// 设置 IsLiked 字段
		for i := range roots {
			if _, liked := likeMap[roots[i].CommentWithUser.ID]; liked {
				roots[i].CommentWithUser.IsLiked = true
			}
		}
	}

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

	// 登录用户则处理点赞
	if userID != "" {
		replyIDs := make([]domain.CommentID, 0, len(replies))
		for i := range replies {
			replyIDs = append(replyIDs, replies[i].CommentWithUser.ID)
		}

		// 整合点赞记录
		likeMap, err := s.cache.GetLikeMap(query.TenantID, userID, replyIDs)
		if err != nil {
			zap.L().Error("获取点赞状态失败", zap.Error(err))
		}

		// 设置 IsLiked 字段
		for i := range replies {
			if _, liked := likeMap[replies[i].CommentWithUser.ID]; liked {
				replies[i].CommentWithUser.IsLiked = true
			}
		}
	}

	return replies, nil
}

func (s *service) ToggleLike(tenantID domain.TenantID, userID domain.UserID, commentID domain.CommentID) error {
	// 去查询当前status
	likeStatus, err := s.cache.GetLikeStatus(tenantID, userID, commentID)
	if err != nil {
		return errors.WithStack(err)
	}

	// toogle状态
	likeStatus.Toogle()

	// 判断当前操作
	if likeStatus.IsLike() {
		// 当前操作为点赞点赞
		if err := s.cache.AddLike(tenantID, userID, commentID); err != nil {
			return errors.WithStack(err)
		}
		if err := s.repo.UpdateLikeCount(tenantID, commentID, likeStatus.IsLike()); err != nil {
			return errors.WithStack(err)
		}
	} else {
		// 当前操作为取消点赞
		if err := s.cache.RemoveLike(tenantID, userID, commentID); err != nil {
			return errors.WithStack(err)
		}
		if err := s.repo.UpdateLikeCount(tenantID, commentID, likeStatus.IsLike()); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
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
			zap.String("tenant_id", string(plate.TenantID)),
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
					zap.String("tenant_id", string(tenantID)),
					zap.String("belong_key", belong.BelongKey),
					zap.String("plate_id", string(belong.ID)),
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
			zap.String("tenant_id", string(config.TenantID)),
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
					zap.String("tenant_id", string(tenantID)),
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
			zap.String("tenant_id", string(config.TenantID)),
			zap.String("plate_id", string(config.Plate.ID)),
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
					zap.String("tenant_id", string(tenantID)),
					zap.String("plate_id", string(plateID)),
				)
			}
		}()
	}

	return config, nil
}
