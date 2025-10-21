package domain

import (
	"database/sql"
)

type TenantRepository interface {
	BeginTx(option ...*sql.TxOptions) (*sql.Tx, error)

	FindByID(id int64) (*Tenant, error)
	FindTenantPlanByID(id int64) (*Plan, error)

	InsertTx(tx *sql.Tx, tenant *Tenant) (*Tenant, error)
	Update(tenant *Tenant) error
	Delete(id int64) error
	List(query *TenantQuery) (*TenantList, error)
}

type TenantCache interface {
}
