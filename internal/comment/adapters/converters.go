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
		PlateID:   comment.Plate.ID,
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
		ID: ormComment.ID,
		Plate: &domain.PlateBelong{
			ID: ormComment.PlateID,
		},
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

func domainPlateToORM(plate *domain.Plate) *orm.CommentPlate {
	if plate == nil {
		return nil
	}

	// 非null项
	ormPlate := &orm.CommentPlate{
		ID:        plate.ID,
		TenantID:  int64(plate.TenantID),
		BelongKey: plate.BelongKey,
		Summary:   plate.Summary,
	}

	// 处理null项

	return ormPlate
}

func ormPlateToDomain(ormPlate *orm.CommentPlate) *domain.Plate {
	if ormPlate == nil {
		return nil
	}

	// 非null项
	plate := &domain.Plate{
		ID:        ormPlate.ID,
		TenantID:  domain.TenantID(ormPlate.TenantID),
		BelongKey: ormPlate.BelongKey,
		Summary:   ormPlate.Summary,
	}

	// 处理null项

	return plate
}

func ormPlatesToDomain(ormPlates []*orm.CommentPlate) []*domain.Plate {
	if len(ormPlates) == 0 {
		return nil
	}

	comments := make([]*domain.Plate, 0, len(ormPlates))
	for i := range ormPlates {
		if ormPlates[i] != nil {
			comments = append(comments, ormPlateToDomain(ormPlates[i]))
		}
	}
	return comments
}

func domainTenantConfigToORM(config *domain.TenantConfig) *orm.CommentTenantConfig {
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

func ormTenantConfigToDomain(ormConfig *orm.CommentTenantConfig) *domain.TenantConfig {
	if ormConfig == nil {
		return nil
	}

	// 非null项
	config := &domain.TenantConfig{
		TenantID:    domain.TenantID(ormConfig.TenantID),
		ClientToken: ormConfig.ClientToken,
		IfAudit:     ormConfig.IfAudit,
		CreatedAt:   ormConfig.CreatedAt,
		UpdatedAt:   ormConfig.UpdatedAt,
	}

	// 处理null项

	return config
}

func domainPlateConfigToORM(config *domain.PlateConfig) *orm.CommentPlateConfig {
	if config == nil {
		return nil
	}

	// 非null项
	ormConfig := &orm.CommentPlateConfig{
		PlateID:  config.Plate.ID,
		TenantID: int64(config.TenantID),
		IfAudit:  config.IfAudit,
	}

	// 处理null项

	return ormConfig
}

func ormPlateConfigToDomain(ormConfig *orm.CommentPlateConfig) *domain.PlateConfig {
	if ormConfig == nil {
		return nil
	}

	// 非null项
	config := &domain.PlateConfig{
		Plate: &domain.PlateBelong{
			ID: ormConfig.PlateID,
		},
		TenantID:  domain.TenantID(ormConfig.TenantID),
		IfAudit:   ormConfig.IfAudit,
		CreatedAt: ormConfig.CreatedAt,
		UpdatedAt: ormConfig.UpdatedAt,
	}

	// 处理null项

	return config
}
