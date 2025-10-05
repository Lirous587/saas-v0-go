package domain

import (
	"database/sql"
	"time"
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

	AssignTenantUserRole(tenantID, userID, roleID int64) error

	ListUsersWithRole(query *UserWithRoleQuery) (*UserWithRoleList, error)
}

type TenantCache interface {
	GenPublicInviteToken(tenantID int64, expireSecond time.Duration) (token string, err error)
	ValidatePublicInviteToken(tenantID int64, value string) error

	GenSecretInviteToken(tenantID int64, expireSecond time.Duration, email string) (token string, err error)
	ValidateSecretInviteToken(tenantID int64, email, value string) error
	DeleteSecretInviteToken(tenantID int64, email string) error
}
