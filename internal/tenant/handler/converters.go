package handler

import (
	"saas/internal/tenant/domain"
)

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

func domainTenantListToResponse(data *domain.TenantList) *TenantListResponse {
	if data == nil {
		return nil
	}

	return &TenantListResponse{
		Total: data.Total,
		List:  domainTenantsToResponse(data.List),
	}
}
