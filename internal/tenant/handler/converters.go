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

func domainUserToResponse(data *domain.User) *UserResponse {
	if data == nil {
		return nil
	}

	return &UserResponse{
		ID:       data.ID,
		Email:    data.Email,
		Nickname: data.Nickname,
	}
}

func domainUsersToResponse(list []*domain.User) []*UserResponse {
	if len(list) == 0 {
		return nil
	}

	ret := make([]*UserResponse, 0, len(list))

	for i := range list {
		if list[i] != nil {
			ret = append(ret, domainUserToResponse(list[i]))
		}
	}
	return ret
}

func domainUserListToResponse(data *domain.UserList) *UserListResponse {
	if data == nil {
		return nil
	}

	return &UserListResponse{
		Total: data.Total,
		List:  domainUsersToResponse(data.List),
	}
}
