package domain

type Role struct {
	ID          int64
	TenantID    int64
	Name        string
	Description string
	IsDefault   bool
}

type RoleQuery struct {
	TenantID int64
	Keyword  string
	Page     int
	PageSize int
}

type RoleList struct {
	Total int64
	List  []*Role
}

func (r *Role) GetDefultManager() int64 {
	return 1
}

func (r *Role) GetDefaultUser() int64 {
	return 2
}
