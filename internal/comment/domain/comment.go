package domain

import (
	"time"
)

type TenantID string

type UserInfo struct {
	ID        string
	NickName  string
	AvatarURL string
	email     string
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

func (cs *CommentStatus) SetApproved() {
	*cs = CommentStatusApproved
}

func (cs *CommentStatus) SetPending() {
	*cs = CommentStatusPending
}

type Comment struct {
	ID        string
	PlateID   string
	UserID    string
	TenantID  TenantID
	ParentID  string
	RootID    string
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
	return c.RootID == "" && c.ParentID == ""
}

func (c *Comment) IsReply() bool {
	return !c.IsRootComment()
}

func (c *Comment) IsReplyRootComment() bool {
	return c.RootID != "" && c.ParentID == ""
}

func (c *Comment) IsReplyParentComment() bool {
	return c.RootID != "" && c.ParentID != ""
}

func (c *Comment) CanReply() bool {
	return c.status == CommentStatusApproved
}

func (c *Comment) IsCommentByAdmin(userID string) bool {
	return c.UserID == userID
}

func (c *Comment) FilterSelf(userIds []string) []string {
	filteredIds := make([]string, 0, 3)
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

type LikeStatus bool

const like LikeStatus = true
const unLike LikeStatus = false

func (l *LikeStatus) Like() {
	*l = like
}

func (l *LikeStatus) UnLike() {
	*l = unLike
}

func (l *LikeStatus) IsLike() bool {
	return *l == like
}

func (l *LikeStatus) Toogle() {
	*l = !*l
}

// -- 评论响应

type CommentWithUser struct {
	ID        string
	User      *UserInfo
	ParentID  string
	RootID    string
	Content   string
	LikeCount int64
	CreatedAt time.Time
	IsLiked   bool
}

type CommentRootsQuery struct {
	TenantID TenantID
	PlateID  string
	LastID   string
	PageSize int
}

type CommentRoot struct {
	CommentWithUser *CommentWithUser
	RepliesCount    int64
}

type CommentRepliesQuery struct {
	TenantID TenantID
	PlateID  string
	RootID   string
	LastID   string
	PageSize int
}

type CommentReply struct {
	CommentWithUser *CommentWithUser
}

// -- 板块

type Plate struct {
	ID         string
	TenantID   TenantID
	BelongKey  string
	RelatedURL string
	Summary    string
}

type PlateBelong struct {
	ID        string
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
