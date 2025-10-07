package handler

import (
    "saas/internal/comment/domain"
)

func domainCommentToResponse(comment *domain.Comment) *CommentResponse {
    if comment == nil {
        return nil
    }

    return &CommentResponse{
        ID:          comment.ID,
        Title:       comment.Title,
        Description: comment.Description,
        CreatedAt:   comment.CreatedAt.Unix(),
        UpdatedAt:   comment.UpdatedAt.Unix(),
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
