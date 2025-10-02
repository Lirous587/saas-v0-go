package adapters

import (
	"github.com/aarondl/null/v8"
	"saas/internal/common/orm"
	"saas/internal/role/domain"
)

func domainRoleToORM(role *domain.Role) *orm.Role {
	if role == nil {
		return nil
	}

	// 非null项
	ormRole := &orm.Role{
		ID:        role.ID,
		Name:      role.Name,
		IsDefault: role.IsDefault,
	}

	// 处理null项
	if role.TenantID != 0 {
		ormRole.TenantID = null.Int64From(role.TenantID)
		ormRole.TenantID.Valid = true
	}

	if role.Description != "" {
		ormRole.Description = null.StringFrom(role.Description)
		ormRole.Description.Valid = true

	}

	return ormRole
}

func ormRoleToDomain(ormRole *orm.Role) *domain.Role {
	if ormRole == nil {
		return nil
	}

	// 非null项
	role := &domain.Role{
		ID:        ormRole.ID,
		Name:      ormRole.Name,
		IsDefault: ormRole.IsDefault,
	}

	// 处理null项
	if ormRole.TenantID.Valid {
		role.TenantID = ormRole.TenantID.Int64
	}
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
