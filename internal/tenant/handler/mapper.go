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

func domainTenantToResponse(tenant *domain.Tenant) *TenantResponse {
	if tenant == nil {
		return nil
	}

	return &TenantResponse{
		ID:          tenant.ID,
		Name:        tenant.Name,
		Description: tenant.Description,
		CreatedAt:   tenant.CreatedAt.Unix(),
		UpdatedAt:   tenant.UpdatedAt.Unix(),
	}
}

func domainTenantsToResponse(tenants []*domain.Tenant) []*TenantResponse {
	if len(tenants) == 0 {
		return nil
	}

	ret := make([]*TenantResponse, 0, len(tenants))

	for _, tenant := range tenants {
		if tenant != nil {
			ret = append(ret, domainTenantToResponse(tenant))
		}
	}
	return ret
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

type UpgradeRequest struct {
	TenantID int64 `uri:"tenant_id" binding:"required"`
	PlanID   int64 `json:"plan_id" binding:"required"`
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

func domainTenantListToResponse(data *domain.TenantList) *TenantListResponse {
	if data == nil {
		return nil
	}

	return &TenantListResponse{
		Total: data.Total,
		List:  domainTenantsToResponse(data.List),
	}
}
