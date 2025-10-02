package domain

type RoleService interface {
	Create(role *Role) (*Role, error)
	Read(id int64) (*Role, error)
	Update(role *Role) (*Role, error)
	Delete(id int64) error
	List(query *RoleQuery) (*RoleList, error)
}
