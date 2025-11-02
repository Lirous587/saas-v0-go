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
		ID:        string(comment.ID),
		PlateID:   string(comment.PlateID),
		UserID:    string(comment.UserID),
		TenantID:  string(comment.TenantID),
		Content:   comment.Content,
		LikeCount: comment.LikeCount,
		CreatedAt: comment.CreatedAt,
		Status:    orm.CommentStatus(comment.Status()),
	}

	// 处理null项
	if comment.RootID != "" {
		ormComment.RootID = null.StringFrom(string(comment.RootID))
		ormComment.RootID.Valid = true
	}
	if comment.ParentID != "" {
		ormComment.ParentID = null.StringFrom(string(comment.ParentID))
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
		ID:        domain.CommentID(ormComment.ID),
		PlateID:   domain.PlateID(ormComment.PlateID),
		UserID:    domain.UserID(ormComment.UserID),
		TenantID:  domain.TenantID(ormComment.TenantID),
		Content:   ormComment.Content,
		LikeCount: ormComment.LikeCount,
		CreatedAt: ormComment.CreatedAt,
	}

	comment.SetPending()

	if ormComment.Status == orm.CommentStatusApproved {
		comment.SetApproved()
	}

	// 处理null项
	if ormComment.RootID.Valid {
		comment.RootID = domain.CommentID(ormComment.RootID.String)
	}
	if ormComment.ParentID.Valid {
		comment.ParentID = domain.CommentID(ormComment.ParentID.String)
	}

	return comment
}

// func ormCommentsToDomain(ormComments []*orm.Comment) []*domain.Comment {
// 	if len(ormComments) == 0 {
// 		return nil
// 	}

// 	comments := make([]*domain.Comment, 0, len(ormComments))
// 	for i := range ormComments {
// 		if ormComments[i] != nil {
// 			comments = append(comments, ormCommentToDomain(ormComments[i]))
// 		}
// 	}
// 	return comments
// }

func ormUserToDomain(ormUser *orm.User) *domain.UserInfo {
	if ormUser == nil {
		return nil
	}

	// 非null项
	user := &domain.UserInfo{
		ID:        domain.UserID(ormUser.ID),
		NickName:  ormUser.Nickname,
		AvatarURL: ormUser.AvatarURL,
	}

	user.SetEmail(ormUser.Email)

	return user
}

func ormUsersToDomain(ormUsers []*orm.User) []*domain.UserInfo {
	if len(ormUsers) == 0 {
		return nil
	}

	users := make([]*domain.UserInfo, 0, len(ormUsers))
	for i := range ormUsers {
		if ormUsers[i] != nil {
			users = append(users, ormUserToDomain(ormUsers[i]))
		}
	}
	return users
}

func domainPlateToORM(plate *domain.Plate) *orm.CommentPlate {
	if plate == nil {
		return nil
	}

	// 非null项
	ormPlate := &orm.CommentPlate{
		ID:         string(plate.ID),
		TenantID:   string(plate.TenantID),
		BelongKey:  plate.BelongKey,
		RelatedURL: plate.RelatedURL,
		Summary:    plate.Summary,
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
		ID:         domain.PlateID(ormPlate.ID),
		TenantID:   domain.TenantID(ormPlate.TenantID),
		BelongKey:  ormPlate.BelongKey,
		RelatedURL: ormPlate.RelatedURL,
		Summary:    ormPlate.Summary,
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
		TenantID: string(config.TenantID),
		IfAudit:  config.IfAudit,
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
		TenantID:  domain.TenantID(ormConfig.TenantID),
		IfAudit:   ormConfig.IfAudit,
		CreatedAt: ormConfig.CreatedAt,
		UpdatedAt: ormConfig.UpdatedAt,
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
		PlateID:  string(config.Plate.ID),
		TenantID: string(config.TenantID),
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
			ID: domain.PlateID(ormConfig.PlateID),
		},
		TenantID:  domain.TenantID(ormConfig.TenantID),
		IfAudit:   ormConfig.IfAudit,
		CreatedAt: ormConfig.CreatedAt,
		UpdatedAt: ormConfig.UpdatedAt,
	}

	// 处理null项

	return config
}
