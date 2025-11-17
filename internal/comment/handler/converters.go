package handler

import (
	"saas/internal/comment/domain"
)

func userInfoToResponse(user *domain.UserInfo) *UserInfo {
	if user == nil {
		return nil
	}
	return &UserInfo{
		ID:       user.ID,
		NickName: user.NickName,
		Avatar:   user.Avatar,
	}
}

func domainCommentRootsToResponse(roots []*domain.CommentRoot) []*CommentRootResponse {
	if len(roots) == 0 {
		return nil
	}

	responses := make([]*CommentRootResponse, 0, len(roots))
	for i := range roots {
		if roots[i] == nil {
			continue
		}

		responses = append(responses, &CommentRootResponse{
			ID:           roots[i].ID,
			User:         userInfoToResponse(roots[i].User),
			Content:      roots[i].Content,
			LikeCount:    roots[i].LikeCount,
			CreatedAt:    roots[i].CreatedAt.Unix(),
			IsLiked:      roots[i].IsLiked,
			RepliesCount: roots[i].RepliesCount,
		})
	}

	return responses
}

func domainCommentRepliesToResponse(replies []*domain.CommentReply) []*CommentReplyResponse {
	if len(replies) == 0 {
		return nil
	}

	responses := make([]*CommentReplyResponse, 0, len(replies))
	for i := range replies {
		if replies[i] == nil {
			continue
		}

		responses = append(responses, &CommentReplyResponse{
			ID:        replies[i].ID,
			ToUser:    userInfoToResponse(replies[i].ToUser),
			User:      userInfoToResponse(replies[i].User),
			ParentID:  replies[i].ParentID,
			RootID:    replies[i].RootID,
			Content:   replies[i].Content,
			LikeCount: replies[i].LikeCount,
			CreatedAt: replies[i].CreatedAt.Unix(),
			IsLiked:   replies[i].IsLiked,
		})
	}

	return responses
}

func domainCommentNoAuditsToResponse(items []*domain.CommentNoAudit) []*CommentNoAuditResponse {
	if len(items) == 0 {
		return nil
	}

	responses := make([]*CommentNoAuditResponse, 0, len(items))
	for _, it := range items {
		if it == nil {
			continue
		}
		responses = append(responses, &CommentNoAuditResponse{
			ID:        it.ID,
			User:      userInfoToResponse(it.User),
			Content:   it.Content,
			CreatedAt: it.CreatedAt.Unix(),
		})
	}

	return responses
}

func domainTenantConfigToResponse(config *domain.TenantConfig) *TenantConfigResponse {
	if config == nil {
		return nil
	}

	return &TenantConfigResponse{
		IfAudit:   config.IfAudit,
		CreatedAt: config.CreatedAt.Unix(),
		UpdatedAt: config.UpdatedAt.Unix(),
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
