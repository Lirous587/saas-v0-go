package domain

import (
	"database/sql"
)

type TenantRepository interface {
	BeginTx(option ...*sql.TxOptions) (*sql.Tx, error)

	FindByID(id int64) (*Tenant, error)

	InsertTx(tx *sql.Tx, tenant *Tenant) (*Tenant, error)
	InsertPlanTx(tx *sql.Tx, tenantID int64, planID int64) error
	InsertUserTx(tx *sql.Tx, tenantID int64, userID int64) error

	Update(tenant *Tenant) (*Tenant, error)
	Delete(id int64) error
	List(query *TenantQuery) (*TenantList, error)
}

type TenantCache interface {
}
