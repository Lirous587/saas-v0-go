package handler

import "saas/internal/comment/domain"

type CommentResponse struct {
	ID        int64                `json:"id"`
	ParentID  int64                `json:"parent_id"`
	RootID    int64                `json:"root_id"`
	Content   string               `json:"content"`
	Status    domain.CommentStatus `json:"status,omitempty"`
	LikeCount int64                `json:"like_count"`
	CreatedAt int64                `json:"created_at"`
	IsLiked   bool                 `json:"is_liked"`
}

type CreateRequest struct {
	ParentID int64 `json:"parent_id"`
	// RootID   int64  `json:"root_id"`
	Content string `json:"content" binding:"required"`
}

type ListRequest struct {
	Page     int    `form:"page,default=1" binding:"min=1"`
	PageSize int    `form:"page_size,default=5" binding:"min=5,max=20"`
	KeyWord  string `form:"keyword" binding:"max=20"`
}

type CommentListResponse struct {
	Total int64              `json:"total"`
	List  []*CommentResponse `json:"list"`
}
