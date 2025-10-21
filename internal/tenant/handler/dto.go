package handler

import (
	"saas/internal/tenant/domain"
)

type TenantResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

type CreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"max=120"`
	PlanID      int64  `json:"plan_id" binding:"required"`
}

type UpdateRequest struct {
	ID          int64  `json:"-" uri:"id" binding:"required"`
	Name        string `json:"id" binding:"required"`
	Description string `json:"description" binding:"max=120"`
}

type DeleteRequest struct {
	ID int64 `json:"-" uri:"id" binding:"required"`
}

type ListRequest struct {
	Page     int    `form:"page,default=1" binding:"min=1"`
	PageSize int    `form:"page_size,default=5" binding:"min=5,max=20"`
	KeyWord  string `form:"keyword" binding:"max=20"`
}

type TenantListResponse struct {
	Total int64             `json:"total"`
	List  []*TenantResponse `json:"list"`
}

type UpgradeRequest struct {
	TenantID int64 `json:"-" uri:"id" binding:"required"`
	PlanID   int64 `json:"plan_id" binding:"required"`
}

type GenInviteTokenRequest struct {
	TenantID     int64 `json:"-" uri:"id" binding:"required"`
	ExpireSecond int64 `json:"expire_second" binding:"required"`
}

type InviteRequest struct {
	TenantID     int64    `json:"-" uri:"id" binding:"required"`
	ExpireSecond int64    `json:"expire_second" binding:"required"`
	Emails       []string `json:"emails" binding:"required,dive,email"`
}

type InviteResponse struct {
	Token string `json:"token"`
}

type EntryRequest struct {
	TenantID  int64                  `json:"-" uri:"id" binding:"required"`
	TokenKind domain.InviteTokenKind `json:"token_kind" form:"token_kind" binding:"required,oneof=public secret"`
	Token     string                 `json:"token" form:"token" binding:"required"`
	Email     string                 `json:"email" form:"number" binding:"required,email"`
}

type ListUserRequest struct {
	TenantID int64  `json:"-" uri:"id" binding:"required"`
	Nickname string `json:"-" form:"nickname"`
	Page     int    `json:"-" form:"page,default=1" binding:"min=1"`
	PageSize int    `json:"-" form:"page_size,default=5" binding:"min=5,max=15"`
}

type UserResponse struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
}

type UserListResponse struct {
	Total int64           `json:"total"`
	List  []*UserResponse `json:"list"`
}
