package domain

import (
	"database/sql"
)

type TenantRepository interface {
	BeginTx(option ...*sql.TxOptions) (*sql.Tx, error)

	GetByID(id string) (*Tenant, error)

	Create(tenant *Tenant) (*Tenant, error)
	Update(tenant *Tenant) error
	Delete(id string) error
	Paging(query *TenantPagingQuery) (*TenantPagination, error)
	ExistSameName(creatorID string, name string) (bool, error)
	IsCreatorHasPlan(creatorID string, planType PlanType) (bool, error)

	GetPlan(id string) (*Plan, error)
}

type TenantCache interface {
}
