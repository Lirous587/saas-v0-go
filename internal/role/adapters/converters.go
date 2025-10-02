package adapters

import (
	"saas/internal/common/orm"
	"saas/internal/role/domain"
	"github.com/aarondl/null/v8"
)

func domainRoleToORM(role *domain.Role) *orm.Role {
	if role == nil {
		return nil
	}

	// 非null项
	ormRole := &orm.Role{
		ID:        		role.ID,
		Title:     		role.Title,
        CreatedAt: 		role.CreatedAt,
        UpdatedAt: 		role.UpdatedAt,
	}

	// 处理null项
	if role.Description != "" {
	 	ormRole.Description = null.StringFrom(role.Description)
	}

	return ormRole
}

func ormRoleToDomain(ormRole *orm.Role) *domain.Role {
	if ormRole == nil {
		return nil
	}

	// 非null项
	role := &domain.Role{
		ID:        		ormRole.ID,
		Title:     		ormRole.Title,
		CreatedAt: 		ormRole.CreatedAt,
		UpdatedAt: 		ormRole.UpdatedAt,
	}

	// 处理null项
	if ormRole.Description.Valid {
 	 	role.Description = ormRole.Description.String
	}

	return role
}

func ormRolesToDomain(ormRoles []*orm.Role) []*domain.Role {
	if len(ormRoles) == 0 {
		return nil
	}

	roles := make([]*domain.Role, 0, len(ormRoles))
	for _, ormRole := range ormRoles {
		if ormRole != nil {
			roles = append(roles, ormRoleToDomain(ormRole))
		}
	}
	return roles
}

