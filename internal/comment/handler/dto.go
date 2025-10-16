package handler

import (
	"saas/internal/comment/domain"
	"saas/internal/common/reskit/codes"
)

type UserInfo struct {
	ID       int64  `json:"id"`
	NickName string `json:"nickname"`
	Avatar   string `json:"avatar,omitempty"`
}

type CreateRequest struct {
	TenantID  domain.TenantID `json:"-" uri:"tenant_id" binding:"required"`
	BelongKey string          `json:"-" uri:"belong_key" binding:"required,max=50"`
	RootID    int64           `json:"root_id"`
	ParentID  int64           `json:"parent_id"`
	Content   string          `json:"content" binding:"required"`
}

// Validate 用于验证root_id和parent_id的组合关系是否正确
func (cq *CreateRequest) Validate() error {
	// 当有parent_id时 必须要有root_id
	if cq.ParentID != 0 {
		if cq.RootID == 0 {
			return codes.ErrCommentIllegalReply
		}
	}
	return nil
}

type DeleteRequest struct {
	TenantID domain.TenantID `json:"-" uri:"tenant_id" binding:"required"`
	ID       int64           `json:"-" uri:"id" binding:"required"`
}

type ListRootsRequest struct {
	TenantID  domain.TenantID `json:"-" uri:"tenant_id" binding:"required"`
	BelongKey string          `json:"-" uri:"belong_key" binding:"required"`
	LastID    int64           `json:"-" form:"last_id" binding:"min=0"`
	PageSize  int             `json:"-" form:"page_size,default=5" binding:"min=5,max=15"`
}

type CommentRootResponse struct {
	ID           int64     `json:"id"`
	User         *UserInfo `json:"user"`
	ParentID     int64     `json:"parent_id,omitempty"`
	RootID       int64     `json:"root_id,omitempty"`
	Content      string    `json:"content"`
	LikeCount    int64     `json:"like_count,omitempty"`
	CreatedAt    int64     `json:"created_at"`
	IsLiked      bool      `json:"is_liked"`
	RepliesCount int64     `json:"replies_count,omitempty"`
}

type ListRepliesRequest struct {
	TenantID  domain.TenantID `json:"-" uri:"tenant_id" binding:"required"`
	BelongKey string          `json:"-" uri:"belong_key" binding:"required"`
	RootID    int64           `json:"-" uri:"root_id" binding:"required"`
	LastID    int64           `json:"-" form:"last_id" binding:"min=0"`
	PageSize  int             `json:"-" form:"page_size,default=5" binding:"min=5,max=15"`
}

type CommentReplyResponse struct {
	ID        int64     `json:"id"`
	User      *UserInfo `json:"user"`
	ParentID  int64     `json:"parent_id,omitempty"`
	RootID    int64     `json:"root_id,omitempty"`
	Content   string    `json:"content"`
	LikeCount int64     `json:"like_count,omitempty"`
	CreatedAt int64     `json:"created_at"`
	IsLiked   bool      `json:"is_liked"`
}

type ToggleLikeRequest struct {
	TenantID domain.TenantID `json:"-" uri:"tenant_id" binding:"required"`
	ID       int64           `json:"-" uri:"id" binding:"required"`
}

type AuditRequest struct {
	TenantID domain.TenantID      `json:"-" uri:"tenant_id" binding:"required"`
	ID       int64                `json:"-" uri:"id" binding:"required"`
	Status   domain.CommentStatus `json:"status" binding:"required"`
}

// --- 评论板块

type PlateResponse struct {
	ID         int64  `json:"id"`
	BelongKey  string `json:"belong_key"`
	RelatedURL string `json:"related_url"`
	Summary    string `json:"summary"`
}

type PlateListResponse struct {
	Total int64            `json:"total"`
	List  []*PlateResponse `json:"list"`
}

type CreatePlateRequest struct {
	TenantID   domain.TenantID `json:"-" uri:"tenant_id" binding:"required"`
	BelongKey  string          `json:"belong_key" binding:"required,max=50"`
	RelatedURL string          `json:"related_url" binding:"required,url,max=255"`
	Summary    string          `json:"summary" binding:"required,max=60"`
}

type UpdatePlateRequest struct {
	TenantID   domain.TenantID `json:"-" uri:"tenant_id" binding:"required"`
	ID         int64           `json:"-" uri:"id" binding:"required"`
	BelongKey  string          `json:"belong_key" binding:"required,max=50"`
	RelatedURL string          `json:"related_url" binding:"required,url,max=255"`
	Summary    string          `json:"summary" binding:"required,max=60"`
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
	TenantID  domain.TenantID `json:"-" uri:"tenant_id" binding:"required"`
	BelongKey string          `json:"-" uri:"belong_key" binding:"required,max=50"`
	IfAudit   *bool           `json:"if_audit" binding:"required"`
}

type GetPlateConfigRequest struct {
	TenantID domain.TenantID `json:"-" uri:"tenant_id" binding:"required"`
	ID       int64           `json:"-" uri:"id" binding:"required"`
}

type PlateConfigResponse struct {
	IfAudit   bool  `json:"if_audit"`
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}
