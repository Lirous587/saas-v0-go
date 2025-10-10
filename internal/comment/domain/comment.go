package domain

import (
	"time"
)

type TenantID int64

type BelongKey string

type UserInfo struct {
	ID       int64  `json:"id"`
	NikeName string `json:"nickname"`
	Avatar   string `json:"avatar,omitempty"`
}

type CommentStatus string

const CommentStatusApprove CommentStatus = ""

type Comment struct {
	ID        int64
	BelongKey BelongKey
	User      *UserInfo
	TenantID  TenantID
	ParentID  int64
	RootID    int64
	Content   string
	Status    CommentStatus
	LikeCount int64
	CreatedAt time.Time
	IsLiked   bool
}

type CommentQuery struct {
	Page     int
	PageSize int
}

type CommentAdvancedQuery struct {
	// Page     int
	// PageSize int
}

type CommentList struct {
	Total int64
	List  []*Comment
}

// CommentTenantConfig 租户全局配置
type CommentTenantConfig struct {
	TenantID    TenantID
	ClientToken string
	IfAudit     bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// CommentConfig belong_key颗粒度的配置 优先级更高
type CommentConfig struct {
	TenantID    TenantID
	BelongKey   BelongKey
	IfAudit     bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
