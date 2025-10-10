package handler

import (
	"saas/internal/comment/domain"
)

func domainCommentToResponse(comment *domain.Comment) *CommentResponse {
	if comment == nil {
		return nil
	}

	return &CommentResponse{
		ID: comment.ID,
		User: &UserInfo{
			ID:       comment.User.ID,
			NikeName: comment.User.NikeName,
			Avatar:   comment.User.Avatar,
		},
		ParentID:  comment.ParentID,
		RootID:    comment.RootID,
		Content:   comment.Content,
		Status:    comment.Status,
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

func domainCommentTenantConfigToResponse(config *domain.CommentTenantConfig) *CommentTenantConfigResponse {
	if config == nil {
		return nil
	}

	return &CommentTenantConfigResponse{
		ClientToken: config.ClientToken,
		IfAudit:     config.IfAudit,
		CreatedAt:   config.CreatedAt.Unix(),
		UpdatedAt:   config.UpdatedAt.Unix(),
	}
}

func domainCommentConfigToResponse(config *domain.CommentConfig) *CommentConfigResponse {
	if config == nil {
		return nil
	}

	return &CommentConfigResponse{
		IfAudit:     config.IfAudit,
		CreatedAt:   config.CreatedAt.Unix(),
		UpdatedAt:   config.UpdatedAt.Unix(),
	}
}
