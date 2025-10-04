package handler

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

type InviteWay string
type InviteRequest struct {
	TenantID int64 `uri:"id" binding:"required"`
	// 1.生成公共令牌 2.生成指定成员令牌(基于邮箱)
	Way InviteWay `json:"way" binding:"required"`
}
