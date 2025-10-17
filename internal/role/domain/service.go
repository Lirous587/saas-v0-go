package domain

type RoleService interface {
	NewRole() *Role

	Create(role *Role) error
	Update(role *Role) error
	Delete(id int64) error
	List() (*RoleList, error)

	// GetUserRoleInTenant 获取指定用户在指定租户下的角色
	GetUserRoleInTenant(userID, tenantID int64) (*Role, error)
}
