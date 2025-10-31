package domain

import (
	"time"
)

type UserInfo struct {
	ID        UserID
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

func (cs *CommentStatus) IsApproved() bool {
	return *cs == CommentStatusApproved
}

type Comment struct {
	ID        CommentID
	PlateID   PlateID
	UserID    UserID
	TenantID  TenantID
	ParentID  CommentID
	RootID    CommentID
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
	return c.RootID.IsZero() && c.ParentID.IsZero()
}

func (c *Comment) IsReply() bool {
	return !c.IsRootComment()
}

func (c *Comment) IsReplyRootComment() bool {
	return !c.RootID.IsZero() && c.ParentID.IsZero()
}

func (c *Comment) IsReplyParentComment() bool {
	return c.RootID.IsZero() && !c.ParentID.IsZero()
}

func (c *Comment) CanReply() bool {
	return c.status == CommentStatusApproved
}

func (c *Comment) IsCommentByAdmin(userID UserID) bool {
	return c.UserID == userID
}

func (c *Comment) FilterSelf(userIDs []UserID) []UserID {
	filteredIDs := make([]UserID, 0, 3)
	for i := range userIDs {
		if userIDs[i] != c.UserID {
			filteredIDs = append(filteredIDs, userIDs[i])
		}
	}

	return filteredIDs
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
	ID        CommentID
	User      *UserInfo
	ParentID  CommentID
	RootID    CommentID
	Content   string
	LikeCount int64
	CreatedAt time.Time
	IsLiked   bool
}

type CommentRootsQuery struct {
	TenantID TenantID
	PlateID  PlateID
	LastID   CommentID
	PageSize int
}

type CommentRoot struct {
	CommentWithUser *CommentWithUser
	RepliesCount    int64
}

type CommentRepliesQuery struct {
	TenantID TenantID
	PlateID  PlateID
	RootID   CommentID
	LastID   CommentID
	PageSize int
}

type CommentReply struct {
	CommentWithUser *CommentWithUser
}

// -- 板块

type Plate struct {
	ID         PlateID
	TenantID   TenantID
	BelongKey  string
	RelatedURL string
	Summary    string
}

type PlateBelong struct {
	ID        PlateID
	BelongKey string
}

type PlateQuery struct {
	TenantID TenantID
	Keyword  string
	Page     int
	PageSize int
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
	TenantID  TenantID
	IfAudit   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// PlateConfig  板块级别的配置 优先级更高
type PlateConfig struct {
	TenantID  TenantID
	Plate     *PlateBelong
	IfAudit   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
