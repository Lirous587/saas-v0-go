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
	Plate    string          `json:"plate" binding:"required"`
	TenantID domain.TenantID `uri:"tenant_id" binding:"required"`
	ParentID int64           `json:"parent_id"`
	Content  string          `json:"content" binding:"required"`
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

// --- 评论板块

type PlateResponse struct {
	ID          int64  `json:"id"`
	Plate       string `json:"plate"`
	Description string `json:"description,omitempty"`
}

type PlateListResponse struct {
	Total int64            `json:"total"`
	List  []*PlateResponse `json:"list"`
}

type CreatePlateRequest struct {
	TenantID    domain.TenantID `json:"-" uri:"tenant_id" binding:"required"`
	Plate       string          `json:"plate" binding:"required,max=50"`
	Description string          `json:"description" binding:"max=60"`
}

type DeletePlateRequest struct {
	TenantID domain.TenantID `json:"-" uri:"tenant_id" binding:"required"`
	ID       int64           `json:"-" uri:"id" binding:"required"`
}

type PlateListRequest struct {
	TenantID domain.TenantID `json:"-" uri:"tenant_id" binding:"required"`
	Page     int             `form:"page,default=1" binding:"min=1"`
	PageSize int             `form:"page_size,default=5" binding:"min=5,max=20"`
	Keyword  string          `form:"keyword" binding:"max=20"`
}

// --- 租户级别配置

type SetTenantConfigRequest struct {
	TenantID domain.TenantID `json:"-" uri:"tenant_id" binding:"required"`
	IfAudit  *bool           `json:"if_audit" binding:"required"`
}

type GetTenantConfigRequest struct {
	TenantID domain.TenantID `json:"-" uri:"tenant_id" binding:"required"`
}

type TenantConfigResponse struct {
	ClientToken string `json:"client_token"`
	IfAudit     bool   `json:"if_audit"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

// --- 板块级别配置

type SetPlateConfigRequest struct {
	TenantID domain.TenantID `json:"-" uri:"tenant_id" binding:"required"`
	Plate    string          `json:"-" uri:"plate" binding:"required"`
	IfAudit  *bool           `json:"if_audit" binding:"required"`
}

type GetPlateConfigRequest struct {
	TenantID domain.TenantID `json:"-" uri:"tenant_id" binding:"required"`
	Plate    string          `json:"-" uri:"plate" binding:"required"`
}

type PlateConfigResponse struct {
	IfAudit   bool  `json:"if_audit"`
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}
