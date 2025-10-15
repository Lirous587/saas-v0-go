package service

import (
	"saas/internal/comment/domain"
	"saas/internal/common/reskit/codes"
	"saas/internal/common/utils"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (s *service) Create(comment *domain.Comment, belongKey string) error {
	// 获取plateBelong 1次sql
	plateBelong, err := s.repo.GetPlateBelongByKey(comment.TenantID, belongKey)
	if err != nil {
		return errors.WithStack(err)
	}

	comment.PlateID = plateBelong.ID

	// 验证评论合法性 1次sql
	if err := s.validateCommentLegitimacy(comment); err != nil {
		return errors.WithStack(err)
	}

	// 查询租户管理员id 一次sql
	admin, err := s.repo.GetDomainAdminByTenant(comment.TenantID)
	if err != nil {
		return errors.WithStack(err)
	}

	// 评论来源:
	// -评论为admin评论
	// 1.无需审计
	// 2.无需通知admin
	// 3.通知viewer(评论通知)
	// -评论为viewer评论
	// 1.检测是否需要审计
	// 2.若要审计 则通知admin(新的评论审核)
	// 3.若无需审计 通知admin viewer(评论通知)
	if comment.IsCommentByAdmin(admin.ID) {
		if err := s.adminCommnet(comment); err != nil {
			return errors.WithStack(err)
		}
	} else {
		if err := s.viewerComment(comment, admin); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (s *service) validateCommentLegitimacy(comment *domain.Comment) error {
	// 请求参数验证确保了root_id和parent_id只有两种组合
	// -无root_id和parent_id
	// 1.仅有root_id（回复root）
	// 2.root_id和parent_id同时存在（回复parent）

	if comment.IsReplyRootComment() {
		rcommnet, err := s.repo.GetByID(comment.TenantID, comment.RootID)
		if err != nil {
			if errors.Is(err, codes.ErrCommentNotFound) {
				return codes.ErrCommentBuildIllegalTree.WithDetail(map[string]any{
					"reason": "当前根评论不存在",
				})
			}
			return errors.WithMessage(err, "获取rcomment失败")
		}

		// 板块一致性检测
		if rcommnet.PlateID != comment.PlateID {
			return codes.ErrCommentRootNotInPlate
		}
		// 是否可回复检测
		if !rcommnet.CanReply() {
			return codes.ErrCommentIllegalReply.WithDetail(map[string]any{
				"reason": "当前评论未审计，无法回复",
			})
		}

	} else if comment.IsReplyParentComment() {
		// 回复评论
		pcomment, err := s.repo.GetByID(comment.TenantID, comment.ParentID)
		if err != nil {
			if errors.Is(err, codes.ErrCommentNotFound) {
				return codes.ErrCommentBuildIllegalTree.WithDetail(map[string]any{
					"reason": "当前父评论不存在",
				})
			}
			return errors.WithMessage(err, "获取pcomment失败")
		}
		// 板块一致性检测 评论层级检测
		if !(pcomment.PlateID == comment.PlateID && pcomment.RootID == comment.RootID) {
			return codes.ErrCommentBuildIllegalTree.WithDetail(map[string]any{
				"reason": "评论层级错误",
			})
		}
		// 是否可回复检测
		if !pcomment.CanReply() {
			return codes.ErrCommentIllegalReply.WithDetail(map[string]any{
				"reason": "当前评论未审计，无法回复",
			})
		}
	}
	return nil
}

// 获取评论配置
func (s *service) getCommentConfig(comment *domain.Comment) (*domain.CommentConfig, error) {
	if comment.TenantID == 0 || comment.PlateID == 0 {
		return nil, codes.ErrCommentIllegalReply
	}

	// 1.先去获取板块配置
	// 2.不存在板块配置则使用租户配置
	plateConfig, err := s.GetPlateConfig(comment.TenantID, comment.PlateID)
	if err == nil {
		zap.L().Debug("使用板块配置")
		return &domain.CommentConfig{
			IfAudit: plateConfig.IfAudit,
		}, nil
	} else if errors.Is(err, codes.ErrCommentPlateConfigNotFound) {
		tenantConfig, err := s.GetTenantConfig(comment.TenantID)
		if err == nil {
			zap.L().Debug("使用租户配置")
			return &domain.CommentConfig{
				IfAudit: tenantConfig.IfAudit,
			}, nil
		} else if errors.Is(err, codes.ErrCommentTenantConfigNotFound) {
			zap.L().Warn("租户配置不存在，使用默认审核")
			return &domain.CommentConfig{
				IfAudit: true,
			}, nil
		}
	}

	// 出现意外错误 记录日志 并且返回默认配置
	zap.L().Error("获取配置失败",
		zap.Int64("tenant_id", int64(comment.TenantID)),
		zap.Int64("comment_id", int64(comment.ID)),
		zap.Int64("plate_id", comment.PlateID),
		zap.Error(err))

	return &domain.CommentConfig{
		IfAudit: true,
	}, nil
}

func (s *service) adminCommnet(comment *domain.Comment) error {
	// 无需审核
	comment.SetApproved()
	// 创建评论
	_, err := s.repo.Create(comment)
	if err != nil {
		return errors.WithStack(err)
	}

	// 异步邮箱通知
	go func() {
		// 如果是创建根级评论 则无需发送邮件通知 不用处理

		// 如果是回复评论
		if comment.IsReply() {
			commnetSource, err := s.getCommentSource(comment)
			if err != nil {
				zap.L().Error("获取评论来源失败", zap.Error(err))
				return
			}

			// 获取root和parent用户id
			uids, err := s.repo.GetUserIdsByRootORParent(comment.TenantID, comment.PlateID, comment.RootID, comment.ParentID)
			if err != nil {
				zap.L().Error("获取root和parent用户id失败", zap.Error(err))
				return
			}

			// 排除自己
			filteredUids := comment.FilterSelf(uids)
			// 去重
			toUserIds := utils.UniqueInt64s(filteredUids)

			// 无toUserIds 无需发送邮件
			if len(toUserIds) == 0 {
				zap.L().Debug("无需发送邮件")
				return
			}

			// 查询所要发送邮件的用户
			toUsers, err := s.repo.GetUserInfosByIds(toUserIds)
			if err != nil {
				zap.L().Error("查询所要发送邮件的用户失败", zap.Error(err))
				return
			}

			// 整合数据 发送邮件
			for _, toUser := range toUsers {
				go func(u *domain.UserInfo) {
					if err := s.sentCommentEmail(commnetSource.commentUser, u.GetEmail(), commnetSource.relatedURL, comment.Content); err != nil {
						zap.L().Error("发送邮件失败", zap.Error(err))
						return
					}
				}(toUser)
			}
		}
	}()

	return nil
}

func (s *service) viewerComment(comment *domain.Comment, admin *domain.UserInfo) error {
	// 查询评论配置
	config, err := s.getCommentConfig(comment)
	if err != nil {
		return errors.WithStack(err)
	}

	// 初始评论状态
	if config.IfAudit {
		comment.SetPending()
	} else {
		comment.SetApproved()
	}
	// 创建评论
	comment, err = s.repo.Create(comment)
	if err != nil {
		return errors.WithStack(err)
	}

	// 异步邮箱通知
	go func() {
		commnetSource, err := s.getCommentSource(comment)
		if err != nil {
			zap.L().Error("获取评论来源失败", zap.Error(err))
			return
		}

		// 评论为根评论时 发邮件给租户
		if comment.IsRootComment() {
			if config.IfAudit {
				if err := s.sentNeedAuditEmail(commnetSource.commentUser, admin.GetEmail(), commnetSource.relatedURL, comment.Content); err != nil {
					zap.L().Error("发送邮件AuditEmail给租户管理员失败", zap.Error(err))
					return
				}
			} else {
				if err := s.sentCommentEmail(commnetSource.commentUser, admin.GetEmail(), commnetSource.relatedURL, comment.Content); err != nil {
					zap.L().Error("发送邮件CommentEmail给租户管理员失败", zap.Error(err))
					return
				}
			}
		} else if comment.IsReply() {
			if config.IfAudit {
				// 发邮件给租户
				if err := s.sentNeedAuditEmail(commnetSource.commentUser, admin.GetEmail(), commnetSource.relatedURL, comment.Content); err != nil {
					zap.L().Error("发邮件给租户失败", zap.Error(err))
					return
				}
			} else {
				// 发邮件给租户以及root和parent
				// 获取root和parent用户id
				uids, err := s.repo.GetUserIdsByRootORParent(comment.TenantID, comment.PlateID, comment.RootID, comment.ParentID)
				if err != nil {
					zap.L().Error("获取root和parent用户id失败", zap.Error(err))
					return
				}

				// 排除自己
				filteredUids := comment.FilterSelf(uids)
				// 去重 并添加管理员(此时可以确保管理员id不为用户id)
				toUserIds := utils.UniqueInt64s(append(filteredUids, admin.ID))

				// 无toUserIds 无需发送邮件
				if len(toUserIds) == 0 {
					zap.L().Debug("无需发送邮件")
					return
				}

				// 查询所要发送邮件的用户
				toUsers, err := s.repo.GetUserInfosByIds(toUserIds)
				if err != nil {
					zap.L().Error("查询所要发送邮件的用户失败", zap.Error(err))
					return
				}

				// 整合数据 发送邮件
				for _, toUser := range toUsers {
					go func(u *domain.UserInfo) {
						if err := s.sentCommentEmail(commnetSource.commentUser, u.GetEmail(), commnetSource.relatedURL, comment.Content); err != nil {
							zap.L().Error("发送邮件失败", zap.Error(err))
							return
						}
					}(toUser)
				}
			}
		}
	}()

	return nil
}

type commnetSource struct {
	commentUser *domain.UserInfo
	relatedURL  string
}

// 获取评论来源 --> 评论者 板块信息
func (s *service) getCommentSource(comment *domain.Comment) (*commnetSource, error) {
	// 获取当前评论用户
	commentUser, err := s.repo.GetUserInfoByID(comment.UserID)
	if err != nil {
		return nil, err
	}
	// 获取 板块的related_url
	relatedURL, err := s.repo.GetPlateRelatedURlByID(comment.TenantID, comment.PlateID)
	if err != nil {
		return nil, err
	}

	return &commnetSource{
		commentUser: commentUser,
		relatedURL:  relatedURL,
	}, nil
}
