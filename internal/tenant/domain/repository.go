package domain

import (
	"database/sql"
)

type TenantRepository interface {
	BeginTx(option ...*sql.TxOptions) (*sql.Tx, error)

	InsertTx(tx *sql.Tx, tenant *Tenant) (*Tenant, error)
	Update(tenant *Tenant) error
	Delete(id int64) error
	Paging(query *TenantPagingQuery) (*TenantPagination, error)
	ExistSameName(creatorID int64, name string) (bool, error)
}

type TenantCache interface {
}
