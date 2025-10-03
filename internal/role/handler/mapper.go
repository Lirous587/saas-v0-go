package handler

import (
	"saas/internal/role/domain"
)

type RoleResponse struct {
	ID          int64  `json:"id"`
	TenantID    int64  `json:"tenant_id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

func domainRoleToResponse(role *domain.Role) *RoleResponse {
	if role == nil {
		return nil
	}

	return &RoleResponse{
		ID:          role.ID,
		TenantID:    role.TenantID,
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

type CreateRequest struct {
	ID          int64  `json:"id" binding:"required"`
	TenantID    int64  `json:"tenant_id"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"max=60"`
	IsDefault   bool   `json:"is_default"`
}

type UpdateRequest struct {
	ID          int64  `json:"id" binding:"required"`
	TenantID    int64  `json:"tenant_id"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"max=60"`
}

type ListRequest struct {
	TenantID int64 `uri:"tenant_id" binding:"required"`
}

type RoleListResponse struct {
	List []*RoleResponse `json:"list"`
}

func domainRoleListToResponse(data *domain.RoleList) *RoleListResponse {
	if data == nil {
		return nil
	}

	return &RoleListResponse{
		List: domainRolesToResponse(data.List),
	}
}
