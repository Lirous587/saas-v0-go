package domain

type TenantService interface {
	Create(tenant *Tenant) error
	Update(tenant *Tenant) error
	Delete(id string) error
	ListByKeyset(query *TenantKeysetQuery) (*TenantKeysetResult, error)
	GetByID(id string) (*Tenant, error)

	CheckName(creatorID string, tenantName string) (bool, error)

	GetPlan(id string) (*Plan, error)
}
