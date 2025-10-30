package domain

type TenantService interface {
	Create(tenant *Tenant) error
	Update(tenant *Tenant) error
	Delete(id string) error
	Paging(query *TenantPagingQuery) (*TenantPagination, error)
	GetByID(id string) (*Tenant, error)

	CheckName(creatorID string, tenantName string) (bool, error)

	GetPlan(id string) (*Plan, error)
}
