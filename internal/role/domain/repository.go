package domain

type RoleRepository interface {
	FindByID(id int64) (*Role, error)

	Create(role *Role) (*Role, error)
	Update(role *Role) error
	Delete(id int64) error
	List() (*RoleList, error)

	FindUserRoleInTenant(userID, tenantID int64) (*Role, error)
}

type RoleCache interface {
	GetUserRoleInTenant(userID, tenantID int64) (*Role, error)
	SetUserRoleInTenant(userID, tenantID int64, role *Role) error
}
