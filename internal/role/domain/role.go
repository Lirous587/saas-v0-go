package domain

import (
	"saas/internal/common/reskit/codes"
)

type Role struct {
	ID          int64
	Name        string
	Description string
	IsDefault   bool
}

type RoleList struct {
	List []*Role
}

func (r *Role) GetTenantadmin() *Role {
	return &Role{
		ID: 1,
	}
}

func (r *Role) GetViewer() *Role {
	return &Role{
		ID: 2,
	}
}

func (r *Role) CheckRoleID(id int64) error {
	switch id {
	case 1: // tenantadmin
		return nil
	case 2: // viewer
		return nil
	default:
		return codes.ErrRoleInvalid
	}
}
