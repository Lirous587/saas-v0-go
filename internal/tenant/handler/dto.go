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
	ID          int64  `uri:"id" binding:"required"`
	Name        string `json:"id" binding:"required"`
	Description string `json:"description" binding:"max=120"`
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
	TenantID int64 `uri:"id" binding:"required"`
	UpgradeRequestBody
}

type UpgradeRequestBody struct {
	PlanID int64 `json:"plan_id" binding:"required"`
}

type GenInviteTokenRequest struct {
	TenantID int64 `uri:"id" binding:"required"`
	GenInviteTokenRequestBody
}

type GenInviteTokenRequestBody struct {
	ExpireSecond int64 `json:"expire_second" binding:"required"`
}

type InviteRequest struct {
	TenantID int64 `uri:"id" binding:"required"`
	InviteRequestBody
}

type InviteRequestBody struct {
	ExpireSecond int64    `json:"expire_second" binding:"required"`
	Emails       []string `json:"emails" binding:"required,dive,email"`
}

type EntryRequest struct {
	TenantID int64 `json:"tenant_id" form:"tenant_id" binding:"required"`
	EntryRequestBody
}

type EntryRequestBody struct {
	TokenKind domain.InviteTokenKind `json:"token_kind" form:"token_kind" binding:"required,oneof=public secret"`
	Token     string                 `json:"token" form:"token" binding:"required"`
	Email     string                 `json:"email" form:"number" binding:"required,email"`
}

type ListUserWithRoleQueryRequest struct {
	TenantID int64 `uri:"id" binding:"required"`
	ListUserWithRoleQueryRequestBody
}

type ListUserWithRoleQueryRequestBody struct {
	RoleID   int64  `form:"role_id"`
	Nickname string `form:"nickname"`

	Page     int `form:"page,default=1" binding:"min=1"`
	PageSize int `form:"page_size,default=5" binding:"min=5,max=20"`
}

type UserResponse struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
}

type RoleResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type UserWithRoleResponse struct {
	User UserResponse `json:"user"`
	Role RoleResponse `json:"role"`
}

type UserWithRoleListResponse struct {
	Total int64                   `json:"total"`
	List  []*UserWithRoleResponse `json:"list"`
}
