package domain

import (
	"database/sql"
)

type TenantRepository interface {
	BeginTx(option ...*sql.TxOptions) (*sql.Tx, error)

	Create(tenant *Tenant) (*Tenant, error)
	Update(tenant *Tenant) error
	Delete(id int64) error
	Paging(query *TenantPagingQuery) (*TenantPagination, error)
	ExistSameName(creatorID int64, name string) (bool, error)
	IsCreatorHasPlan(creatorID int64, planType PlanType) (bool, error)

	GetPlan(id int64) (*Plan, error)
}

type TenantCache interface {
}
