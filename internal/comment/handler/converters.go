package handler

import (
	"saas/internal/comment/domain"
)

func domainCommentToResponse(comment *domain.Comment) *CommentResponse {
	if comment == nil {
		return nil
	}

	return &CommentResponse{
		ID:        comment.ID,
		UserID:    comment.UserID,
		ParentID:  comment.ParentID,
		RootID:    comment.RootID,
		Content:   comment.Content,
		Status:    comment.GetStatus(),
		LikeCount: comment.LikeCount,
		CreatedAt: comment.CreatedAt.Unix(),
		IsLiked:   comment.IsLiked,
	}
}

func domainCommentsToResponse(comments []*domain.Comment) []*CommentResponse {
	if len(comments) == 0 {
		return nil
	}

	ret := make([]*CommentResponse, 0, len(comments))

	for _, comment := range comments {
		if comment != nil {
			ret = append(ret, domainCommentToResponse(comment))
		}
	}
	return ret
}

func domainCommentListToResponse(data *domain.CommentList) *CommentListResponse {
	if data == nil {
		return nil
	}

	return &CommentListResponse{
		Total: data.Total,
		List:  domainCommentsToResponse(data.List),
	}
}

func domainTenantConfigToResponse(config *domain.TenantConfig) *TenantConfigResponse {
	if config == nil {
		return nil
	}

	return &TenantConfigResponse{
		ClientToken: config.ClientToken,
		IfAudit:     config.IfAudit,
		CreatedAt:   config.CreatedAt.Unix(),
		UpdatedAt:   config.UpdatedAt.Unix(),
	}
}

func domainPlateConfigToResponse(config *domain.PlateConfig) *PlateConfigResponse {
	if config == nil {
		return nil
	}

	return &PlateConfigResponse{
		IfAudit:   config.IfAudit,
		CreatedAt: config.CreatedAt.Unix(),
		UpdatedAt: config.UpdatedAt.Unix(),
	}
}

func domainPlateToResponse(plate *domain.Plate) *PlateResponse {
	if plate == nil {
		return nil
	}

	return &PlateResponse{
		ID:         plate.ID,
		BelongKey:  plate.BelongKey,
		Summary:    plate.Summary,
		RelatedURL: plate.RelatedURL,
	}
}

func domainPlatesToResponse(plates []*domain.Plate) []*PlateResponse {
	if len(plates) == 0 {
		return nil
	}

	ret := make([]*PlateResponse, 0, len(plates))

	for i := range plates {
		if plates[i] != nil {
			ret = append(ret, domainPlateToResponse(plates[i]))
		}
	}
	return ret
}

func domainPlateListToResponse(data *domain.PlateList) *PlateListResponse {
	if data == nil {
		return nil
	}

	return &PlateListResponse{
		Total: data.Total,
		List:  domainPlatesToResponse(data.List),
	}
}
