package domain

import (
	"time"
)

type TenantID int64

type UserInfo struct {
	ID       int64
	NickName string
	Avatar   string
	email    string
}

func (u *UserInfo) SetEmail(email string) {
	u.email = email
}

func (u *UserInfo) GetEmail() string {
	return u.email
}

type CommentStatus string

const CommentStatusApprove CommentStatus = ""

type Comment struct {
	ID      int64
	PlateID int64
	UserID  int64
	// Plate     *PlateBelong
	// User      *UserInfo
	TenantID  TenantID
	ParentID  int64
	RootID    int64
	Content   string
	Status    CommentStatus
	LikeCount int64
	CreatedAt time.Time
	IsLiked   bool
}

// IsTopLevelComment 判断是否为顶级评论（root 评论，无父评论）
func (c *Comment) IsTopLevelComment() bool {
	return c.RootID != 0 && c.ParentID == 0 // 确保是 root 且无 parent
}

// IsReplyComment 判断是否为回复评论（有父评论）
func (c *Comment) IsReplyParentComment() bool {
	return c.ParentID != 0
}

func (c *Comment) IsCommentByAdmin(userID int64) bool {
	return c.UserID == userID
}

// IsReply 存在root或parent
func (c *Comment) IsReply() bool {
	return c.ParentID != 0 || c.RootID != 0
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
