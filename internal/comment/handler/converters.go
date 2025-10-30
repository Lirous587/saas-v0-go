package handler

import (
	"saas/internal/comment/domain"
)

func domainCommentRootsToResponse(roots []*domain.CommentRoot) []*CommentRootResponse {
	if len(roots) == 0 {
		return nil
	}

	responses := make([]*CommentRootResponse, 0, len(roots))
	for i := range roots {
		if roots[i] == nil || roots[i].CommentWithUser == nil {
			continue
		}
		CommentWithUser := roots[i].CommentWithUser
		userInfo := &UserInfo{
			ID:        CommentWithUser.User.ID,
			NickName:  CommentWithUser.User.NickName,
			AvatarURL: CommentWithUser.User.AvatarURL,
		}
		responses = append(responses, &CommentRootResponse{
			ID:           CommentWithUser.ID,
			User:         userInfo,
			ParentID:     CommentWithUser.ParentID,
			RootID:       CommentWithUser.RootID,
			Content:      CommentWithUser.Content,
			LikeCount:    CommentWithUser.LikeCount,
			CreatedAt:    CommentWithUser.CreatedAt.Unix(),
			IsLiked:      CommentWithUser.IsLiked,
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
		if replies[i] == nil || replies[i].CommentWithUser == nil {
			continue
		}
		CommentWithUser := replies[i].CommentWithUser
		userInfo := &UserInfo{
			ID:        CommentWithUser.User.ID,
			NickName:  CommentWithUser.User.NickName,
			AvatarURL: CommentWithUser.User.AvatarURL,
		}
		responses = append(responses, &CommentReplyResponse{
			ID:        CommentWithUser.ID,
			User:      userInfo,
			ParentID:  CommentWithUser.ParentID,
			RootID:    CommentWithUser.RootID,
			Content:   CommentWithUser.Content,
			LikeCount: CommentWithUser.LikeCount,
			CreatedAt: CommentWithUser.CreatedAt.Unix(),
			IsLiked:   CommentWithUser.IsLiked,
		})
	}

	return responses
}

func domainTenantConfigToResponse(config *domain.TenantConfig) *TenantConfigResponse {
	if config == nil {
		return nil
	}

	return &TenantConfigResponse{
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
