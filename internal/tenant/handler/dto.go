package handler

type TenantResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

type CreateRequest struct {
	Name        string `json:"name" binding:"required,max=20"`
	Description string `json:"description" binding:"max=120"`
	PlanID      int64  `json:"plan_id" binding:"required"`
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
	PrevCursor string            `json:"prev_cursor"`
	NextCursor string            `json:"next_cursor"`
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
