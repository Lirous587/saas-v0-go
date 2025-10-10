package adapters

import (
	"saas/internal/comment/domain"
	"saas/internal/common/orm"

	"github.com/aarondl/null/v8"
)

func domainCommentToORM(comment *domain.Comment) *orm.Comment {
	if comment == nil {
		return nil
	}

	// 非null项
	ormComment := &orm.Comment{
		ID:        comment.ID,
		BelongKey: string(comment.BelongKey),
		UserID:    comment.User.ID,
		TenantID:  int64(comment.TenantID),
		LikeCount: comment.LikeCount,
		CreatedAt: comment.CreatedAt,
	}

	// 处理null项
	if comment.RootID != 0 {
		ormComment.RootID = null.Int64From(comment.RootID)
		ormComment.RootID.Valid = true
	}
	if comment.ParentID != 0 {
		ormComment.ParentID = null.Int64From(comment.ParentID)
		ormComment.ParentID.Valid = true
	}

	return ormComment
}

func ormCommentToDomain(ormComment *orm.Comment) *domain.Comment {
	if ormComment == nil {
		return nil
	}

	// 非null项
	comment := &domain.Comment{
		ID:        ormComment.ID,
		BelongKey: domain.BelongKey(ormComment.BelongKey),
		User: &domain.UserInfo{
			ID: ormComment.UserID,
			// Avatar: "",
		},
		TenantID:  domain.TenantID(ormComment.TenantID),
		LikeCount: ormComment.LikeCount,
		CreatedAt: ormComment.CreatedAt,
	}

	// 处理null项
	if ormComment.RootID.Valid {
		comment.RootID = ormComment.RootID.Int64
	}
	if ormComment.ParentID.Valid {
		comment.ParentID = ormComment.ParentID.Int64
	}

	return comment
}

func ormCommentsToDomain(ormComments []*orm.Comment) []*domain.Comment {
	if len(ormComments) == 0 {
		return nil
	}

	comments := make([]*domain.Comment, 0, len(ormComments))
	for _, ormComment := range ormComments {
		if ormComment != nil {
			comments = append(comments, ormCommentToDomain(ormComment))
		}
	}
	return comments
}

func domainCommentTenantConfigToORM(config *domain.CommentTenantConfig) *orm.CommentTenantConfig {
	if config == nil {
		return nil
	}

	// 非null项
	ormConfig := &orm.CommentTenantConfig{
		TenantID:    int64(config.TenantID),
		ClientToken: config.ClientToken,
		IfAudit:     config.IfAudit,
	}

	// 处理null项

	return ormConfig
}

func ormCommentTenantConfigToDomain(ormConfig *orm.CommentTenantConfig) *domain.CommentTenantConfig {
	if ormConfig == nil {
		return nil
	}

	// 非null项
	config := &domain.CommentTenantConfig{
		TenantID:    domain.TenantID(ormConfig.TenantID),
		ClientToken: ormConfig.ClientToken,
		IfAudit:     ormConfig.IfAudit,
		CreatedAt:   ormConfig.CreatedAt,
		UpdatedAt:   ormConfig.UpdatedAt,
	}

	// 处理null项

	return config
}

func domainCommentConfigToORM(config *domain.CommentConfig) *orm.CommentConfig {
	if config == nil {
		return nil
	}

	// 非null项
	ormConfig := &orm.CommentConfig{
		BelongKey:   string(config.BelongKey),
		TenantID:    int64(config.TenantID),
		ClientToken: config.ClientToken,
		IfAudit:     config.IfAudit,
	}

	// 处理null项

	return ormConfig
}

func ormCommentConfigToDomain(ormConfig *orm.CommentConfig) *domain.CommentConfig {
	if ormConfig == nil {
		return nil
	}

	// 非null项
	config := &domain.CommentConfig{
		BelongKey:   domain.BelongKey(ormConfig.BelongKey),
		TenantID:    domain.TenantID(ormConfig.TenantID),
		ClientToken: ormConfig.ClientToken,
		IfAudit:     ormConfig.IfAudit,
		CreatedAt:   ormConfig.CreatedAt,
		UpdatedAt:   ormConfig.UpdatedAt,
	}

	// 处理null项

	return config
}
