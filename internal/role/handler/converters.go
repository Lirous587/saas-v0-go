package handler

import (
	"saas/internal/role/domain"
)

func domainRoleToResponse(role *domain.Role) *RoleResponse {
	if role == nil {
		return nil
	}

	return &RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
	}
}

func domainRolesToResponse(roles []*domain.Role) []*RoleResponse {
	if len(roles) == 0 {
		return nil
	}

	ret := make([]*RoleResponse, 0, len(roles))

	for _, role := range roles {
		if role != nil {
			ret = append(ret, domainRoleToResponse(role))
		}
	}
	return ret
}

func domainRoleListToResponse(data *domain.RoleList) *RoleListResponse {
	if data == nil {
		return nil
	}

	return &RoleListResponse{
		List: domainRolesToResponse(data.List),
	}
}
