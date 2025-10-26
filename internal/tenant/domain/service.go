package domain

type TenantService interface {
	Create(tenant *Tenant) error
	Update(tenant *Tenant) error
	Delete(id int64) error
	Paging(query *TenantPagingQuery) (*TenantPagination, error)

	CheckName(creatorID int64, tenantName string) (bool, error)

	GetPlan(id int64) (*Plan, error)
}
