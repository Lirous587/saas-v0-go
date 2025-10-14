package domain

import (
	"time"
)

type TenantID int64

type UserInfo struct {
	ID       int64
	NickName string
	email    string
}

func (u *UserInfo) SetEmail(email string) {
	u.email = email
}

func (u *UserInfo) GetEmail() string {
	return u.email
}

type CommentStatus string

const CommentStatusApproved CommentStatus = "approved"
const CommentStatusPending CommentStatus = "pending"

type Comment struct {
	ID        int64
	PlateID   int64
	UserID    int64
	TenantID  TenantID
	ParentID  int64
	RootID    int64
	Content   string
	status    CommentStatus
	LikeCount int64
	CreatedAt time.Time
	IsLiked   bool
}

func (c *Comment) GetStatus() CommentStatus {
	return c.status
}

func (c *Comment) SetApproved() {
	c.status = CommentStatusApproved
}

func (c *Comment) SetPending() {
	c.status = CommentStatusPending
}

func (c *Comment) IsApproved() bool {
	return c.status == CommentStatusApproved
}

func (c *Comment) IsRootComment() bool {
	return c.RootID == 0 && c.ParentID == 0
}

func (c *Comment) IsReply() bool {
	return !c.IsRootComment()
}

func (c *Comment) IsReplyRootComment() bool {
	return c.RootID != 0 && c.ParentID == 0
}

func (c *Comment) IsReplyParentComment() bool {
	return c.RootID != 0 && c.ParentID != 0
}

func (c *Comment) CanReply() bool {
	return c.status == CommentStatusApproved
}

func (c *Comment) IsCommentByAdmin(userID int64) bool {
	return c.UserID == userID
}

func (c *Comment) FilterSelf(userIds []int64) []int64 {
	filteredIds := make([]int64, 0, 3)
	for i := range userIds {
		if userIds[i] != c.UserID {
			filteredIds = append(filteredIds, userIds[i])
		}
	}

	return filteredIds
}

func (c *Comment) CanAudit() bool {
	return c.status == CommentStatusPending
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

type Plate struct {
	ID         int64
	TenantID   TenantID
	BelongKey  string
	RelatedURL string
	Summary    string
}

type PlateBelong struct {
	ID        int64
	BelongKey string
}

type PlateQuery struct {
	TenantID TenantID
	Page     int
	PageSize int
	Keyword  string
}

type PlateList struct {
	Total int64
	List  []*Plate
}

type CommentConfig struct {
	IfAudit bool
}

// TenantConfig 基于租户的全局配置
type TenantConfig struct {
	TenantID    TenantID
	ClientToken string
	IfAudit     bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// PlateConfig  板块级别的配置 优先级更高
type PlateConfig struct {
	TenantID  TenantID
	Plate     *PlateBelong
	IfAudit   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
