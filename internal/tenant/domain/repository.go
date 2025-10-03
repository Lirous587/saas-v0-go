package domain

import (
	"database/sql"
)

type TenantRepository interface {
	BeginTx(option ...*sql.TxOptions) (*sql.Tx, error)

	FindByID(id int64) (*Tenant, error)
	FindTenantPlanByID(id int64) (*Plan, error)

	InsertTx(tx *sql.Tx, tenant *Tenant) (*Tenant, error)

	Update(tenant *Tenant) (*Tenant, error)
	Delete(id int64) error
	List(query *TenantQuery) (*TenantList, error)

	AssignTenantUserRoleTx(tx *sql.Tx, tenantID, userID, roleID int64) error
}

type TenantCache interface {
}
