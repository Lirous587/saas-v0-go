package handler

import (
	"saas/internal/comment/domain"
	"saas/internal/common/reskit/codes"
)

type UserInfo struct {
	ID        domain.UserID `json:"id"`
	NickName  string        `json:"nickname"`
	AvatarURL string        `json:"avatar_url,omitempty"`
}

type CreateRequest struct {
	TenantID  domain.TenantID  `json:"-" uri:"tenant_id" binding:"required,uuid"`
	BelongKey string           `json:"-" uri:"belong_key" binding:"required,max=50"`
	RootID    domain.CommentID `json:"root_id" binding:"omitempty,uuid"`
	ParentID  domain.CommentID `json:"parent_id" binding:"omitempty,uuid"`
	Content   string           `json:"content" binding:"required"`
}

// Validate 用于验证root_id和parent_id的组合关系是否正确
func (cq *CreateRequest) Validate() error {
	// 当有parent_id时 必须要有root_id
	if !cq.ParentID.IsZero() {
		if cq.RootID.IsZero() {
			return codes.ErrCommentIllegalReply
		}
	}
	return nil
}

type DeleteRequest struct {
	TenantID domain.TenantID  `json:"-" uri:"tenant_id" binding:"required,uuid"`
	ID       domain.CommentID `json:"-" uri:"id" binding:"required,uuid"`
}

type ListRootsRequest struct {
	TenantID  domain.TenantID  `json:"-" uri:"tenant_id" binding:"required,uuid"`
	BelongKey string           `json:"-" uri:"belong_key" binding:"required"`
	LastID    domain.CommentID `json:"-" form:"last_id" binding:"omitempty,uuid"`
	PageSize  int              `json:"-" form:"page_size,default=5" binding:"min=5,max=15"`
}

type CommentRootResponse struct {
	ID           domain.CommentID `json:"id"`
	User         *UserInfo        `json:"user"`
	ParentID     domain.CommentID `json:"parent_id,omitempty"`
	RootID       domain.CommentID `json:"root_id,omitempty"`
	Content      string           `json:"content"`
	LikeCount    int64            `json:"like_count,omitempty"`
	CreatedAt    int64            `json:"created_at"`
	IsLiked      bool             `json:"is_liked,omitempty"`
	RepliesCount int64            `json:"replies_count,omitempty"`
}

type ListRepliesRequest struct {
	TenantID  domain.TenantID  `json:"-" uri:"tenant_id" binding:"required,uuid"`
	BelongKey string           `json:"-" uri:"belong_key" binding:"required"`
	RootID    domain.CommentID `json:"-" uri:"root_id" binding:"required,uuid"`
	LastID    domain.CommentID `json:"-" form:"last_id" binding:"omitempty,uuid"`
	PageSize  int              `json:"-" form:"page_size,default=5" binding:"min=5,max=15"`
}

type CommentReplyResponse struct {
	ID        domain.CommentID `json:"id"`
	ToUser    *UserInfo        `json:"to_user,omitempty"`
	User      *UserInfo        `json:"user"`
	ParentID  domain.CommentID `json:"parent_id,omitempty"`
	RootID    domain.CommentID `json:"root_id,omitempty"`
	Content   string           `json:"content"`
	LikeCount int64            `json:"like_count,omitempty"`
	CreatedAt int64            `json:"created_at"`
	IsLiked   bool             `json:"is_liked,omitempty"`
}

type ToggleLikeRequest struct {
	TenantID domain.TenantID  `json:"-" uri:"tenant_id" binding:"required,uuid"`
	ID       domain.CommentID `json:"-" uri:"id" binding:"required,uuid"`
}

type ListNoAuditRequest struct {
	TenantID  domain.TenantID `json:"-" uri:"tenant_id" binding:"required,uuid"`
	BelongKey string          `json:"-" form:"belong_key" binding:"required"`
	Keyword   string          `form:"keyword"`
	PageSize  int             `json:"-" form:"page_size,default=5" binding:"min=5,max=15"`
}

type CommentNoAuditResponse struct {
	ID        domain.CommentID `json:"id"`
	User      *UserInfo        `json:"user"`
	Content   string           `json:"content"`
	CreatedAt int64            `json:"created_at"`
}

type AuditAction string

func (ac *AuditAction) isAccept() bool {
	return *ac == auditAccept
}

const auditAccept AuditAction = "accept"
const auditReject AuditAction = "reject"

type AuditRequest struct {
	TenantID domain.TenantID  `json:"-" uri:"tenant_id" binding:"required,uuid"`
	ID       domain.CommentID `json:"-" uri:"id" binding:"required,uuid"`
	Action   AuditAction      `json:"action" binding:"required"`
}

// --- 评论板块

type PlateResponse struct {
	ID         domain.PlateID `json:"id"`
	BelongKey  string         `json:"belong_key"`
	RelatedURL string         `json:"related_url"`
	Summary    string         `json:"summary"`
}

type PlateListResponse struct {
	Total int64            `json:"total"`
	List  []*PlateResponse `json:"list"`
}

type CreatePlateRequest struct {
	TenantID   domain.TenantID `json:"-" uri:"tenant_id" binding:"required,uuid"`
	BelongKey  string          `json:"belong_key" binding:"required,max=50"`
	RelatedURL string          `json:"related_url" binding:"required,url,max=255"`
	Summary    string          `json:"summary" binding:"required,max=500"`
}

type UpdatePlateRequest struct {
	TenantID   domain.TenantID `json:"-" uri:"tenant_id" binding:"required,uuid"`
	ID         domain.PlateID  `json:"-" uri:"id" binding:"required,uuid"`
	BelongKey  string          `json:"belong_key" binding:"required,max=50"`
	RelatedURL string          `json:"related_url" binding:"required,url,max=255"`
	Summary    string          `json:"summary" binding:"required,max=60"`
}

type DeletePlateRequest struct {
	TenantID domain.TenantID `json:"-" uri:"tenant_id" binding:"required,uuid"`
	ID       domain.PlateID  `json:"-" uri:"id" binding:"required,uuid"`
}

type PlateListRequest struct {
	TenantID domain.TenantID `json:"-" uri:"tenant_id" binding:"required,uuid"`
	Page     int             `form:"page,default=1" binding:"min=1"`
	PageSize int             `form:"page_size,default=5" binding:"min=5,max=20"`
	Keyword  string          `form:"keyword"`
}

// --- 租户级别配置

type SetTenantConfigRequest struct {
	TenantID domain.TenantID `json:"-" uri:"tenant_id" binding:"required,uuid"`
	IfAudit  *bool           `json:"if_audit" binding:"required"`
}

type GetTenantConfigRequest struct {
	TenantID domain.TenantID `json:"-" uri:"tenant_id" binding:"required,uuid"`
}

type TenantConfigResponse struct {
	IfAudit   bool  `json:"if_audit"`
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

// --- 板块级别配置

type SetPlateConfigRequest struct {
	TenantID  domain.TenantID `json:"-" uri:"tenant_id" binding:"required,uuid"`
	BelongKey string          `json:"-" uri:"belong_key" binding:"required,max=50"`
	IfAudit   *bool           `json:"if_audit" binding:"required"`
}

type GetPlateConfigRequest struct {
	TenantID domain.TenantID `json:"-" uri:"tenant_id" binding:"required"`
	ID       domain.PlateID  `json:"-" uri:"id" binding:"required"`
}

type PlateConfigResponse struct {
	IfAudit   bool  `json:"if_audit"`
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

type PlateCheckBelongKeyRequest struct {
	TenantID  domain.TenantID `json:"-" uri:"tenant_id" binding:"required,uuid"`
	BelongKey string          `json:"-" form:"belong_key" binding:"required"`
}
