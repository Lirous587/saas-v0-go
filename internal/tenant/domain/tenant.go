package domain

import (
	"time"
)

type Tenant struct {
	ID          int64
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CreatorID   int64
}

type Plan struct {
	ID   int64
	Name string
}

type TenantQuery struct {
	Keyword  string
	Page     int
	PageSize int
}

type TenantList struct {
	Total int64
	List  []*Tenant
}

// InviteTokenKind 1.生成公共令牌 2.生成指定成员令牌(基于邮箱)
type InviteTokenKind string

const publicWay InviteTokenKind = "public"
const secretWay InviteTokenKind = "secret"

func (way InviteTokenKind) IsPublic() bool {
	return way == publicWay
}

type GenInviteTokenPayload struct {
	TenantID     int64
	ExpireSecond int64 `json:"expire_second" binding:"required"`
}

type InvitePayload struct {
	TenantID     int64
	ExpireSecond int64 `json:"expire_second" binding:"required"`
	Emails       []string
}

type EnterPayload struct {
	TenantID  int64
	TokenKind InviteTokenKind
	Token     string
	Email     string
}

type UserQuery struct {
	TenantID int64
	Nickname string
	Page     int
	PageSize int
}

type User struct {
	ID       int64
	Email    string
	Nickname string
}

type Role struct {
	ID   int64
	Name string
}

type UserList struct {
	Total int64
	List  []*User
}
