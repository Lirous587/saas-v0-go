package handler

import "saas/internal/tenant/domain"

type TenantResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

type CreateRequest struct {
	Name         string          `json:"name" binding:"required,max=20"`
	Description  string          `json:"description" binding:"max=120"`
	PlanType     domain.PlanType `json:"plan_type" binding:"required,oneof=free caring professional"`
	BillingCycle string          `json:"billing_cycle" binding:"required,oneof=monthly yearly lifetime"`
}

type UpdateRequest struct {
	ID          int64  `json:"-" uri:"id" binding:"required"`
	Name        string `json:"id" binding:"required,max=20"`
	Description string `json:"description" binding:"max=120"`
}

type DeleteRequest struct {
	ID int64 `json:"-" uri:"id" binding:"required"`
}

type PagingRequest struct {
	PageSize   int    `form:"page_size,default=5" binding:"min=5,max=20"`
	PrevCursor string `form:"prev_cursor"`
	NextCursor string `form:"next_cursor"`
	KeyWord    string `form:"keyword" binding:"max=20"`
}

type PagingResponse struct {
	Items      []*TenantResponse `json:"items"`
	PrevCursor string            `json:"prev_cursor,omitempty"`
	NextCursor string            `json:"next_cursor,omitempty"`
	HasPrev    bool              `json:"has_prev"`
	HasNext    bool              `json:"has_next"`
}

type CheckNameRequest struct {
	Name string `form:"name"  binding:"required,max=20"`
}

type UpgradeRequest struct {
	TenantID int64 `json:"-" uri:"id" binding:"required"`
	PlanID   int64 `json:"plan_id" binding:"required"`
}

type GetPlanRequest struct {
	ID int64 `json:"-" uri:"id" binding:"required"`
}

type PlanResponse struct {
	TenantID     int64                   `json:"tenant_id"`
	PlanType     domain.PlanType         `json:"plan_type"`
	StartTime    int64                   `json:"start_time"`
	EndTime      int64                   `json:"end_time"`
	Status       domain.PlanStatus       `json:"status"`
	BillingCycle domain.PlanBillingCycle `json:"billing_cycle"`
	CanUpgrade   bool                    `json:"can_upgrade"`
}
