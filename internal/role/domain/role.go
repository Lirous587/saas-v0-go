package domain

type Role struct {
	ID          int64
	Name        string
	Description string
	IsDefault   bool
}

type RoleList struct {
	List []*Role
}

func (r *Role) GetDefultSuperadmin() *Role {
	return &Role{
		ID: 1,
	}
}

func (r *Role) GetDefaultViewer() *Role {
	return &Role{
		ID: 2,
	}
}
