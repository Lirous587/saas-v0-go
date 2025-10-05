package handler

type RoleResponse struct {
	ID          int64  `json:"id"`
	TenantID    int64  `json:"tenant_id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type CreateRequest struct {
	ID          int64  `json:"id" binding:"required"`
	TenantID    int64  `json:"tenant_id"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"max=60"`
}

type UpdateRequest struct {
	ID          int64  `json:"id" binding:"required"`
	TenantID    int64  `json:"tenant_id"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"max=60"`
}

type RoleListResponse struct {
	List []*RoleResponse `json:"list"`
}
