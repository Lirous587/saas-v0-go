package handler

import "saas/internal/comment/domain"

type UserInfo struct {
	ID       int64  `json:"id"`
	NikeName string `json:"nickname"`
	Avatar   string `json:"avatar,omitempty"`
}

type CommentResponse struct {
	ID        int64                `json:"id"`
	User      *UserInfo            `json:"user"`
	ParentID  int64                `json:"parent_id"`
	RootID    int64                `json:"root_id"`
	Content   string               `json:"content"`
	Status    domain.CommentStatus `json:"status,omitempty"`
	LikeCount int64                `json:"like_count"`
	CreatedAt int64                `json:"created_at"`
	IsLiked   bool                 `json:"is_liked"`
}

type CreateRequest struct {
	BelongKey domain.BelongKey `json:"belong_key" binding:"required"`
	TenantID  domain.TenantID  `uri:"tenant_id" binding:"required"`
	ParentID  int64            `json:"parent_id"`
	Content   string           `json:"content" binding:"required"`
}

type ListRequest struct {
	Page     int    `form:"page,default=1" binding:"min=1"`
	PageSize int    `form:"page_size,default=5" binding:"min=5,max=20"`
	KeyWord  string `form:"keyword" binding:"max=20"`
}

type AdvancedListRequest struct {
}

type CommentListResponse struct {
	Total int64              `json:"total"`
	List  []*CommentResponse `json:"list"`
}

type CommentTenantConfigResponse struct {
	ClientToken string `json:"client_token"`
	IfAudit     bool   `json:"if_audit"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

type SetCommentTenantConfigRequest struct {
	TenantID domain.TenantID `json:"-" uri:"tenant_id" binding:"required"`
	IfAudit  *bool           `json:"if_audit" binding:"required"`
}

type GetCommentTenantConfigRequest struct {
	TenantID domain.TenantID `json:"-" uri:"tenant_id" binding:"required"`
}

type CommentConfigResponse struct {
	IfAudit     bool   `json:"if_audit"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

type SetCommentConfigRequest struct {
	TenantID  domain.TenantID  `json:"-" uri:"tenant_id" binding:"required"`
	BelongKey domain.BelongKey `json:"-" uri:"belong_key" binding:"required"`
	IfAudit   *bool            `json:"if_audit" binding:"required"`
}

type GetCommentConfigRequest struct {
	TenantID  domain.TenantID  `json:"-" uri:"tenant_id" binding:"required"`
	BelongKey domain.BelongKey `json:"-" uri:"belong_key" binding:"required"`
}
