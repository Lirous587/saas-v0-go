package domain


type RoleRepository interface {
	FindByID(id int64) (*Role, error)

	Create(role *Role) (*Role, error)
	Update(role *Role) (*Role, error)
	Delete(id int64) error
	List(query *RoleQuery) (*RoleList, error)
}

type RoleCache interface {

}
