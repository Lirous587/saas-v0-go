package adapters

import (
	"context"
	"database/sql"
	"fmt"
	"saas/internal/common/orm"
	"saas/internal/common/reskit/codes"
	"saas/internal/tenant/domain"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/pkg/errors"
)

type TenantPSQLRepository struct {
}

func NewTenantPSQLRepository() domain.TenantRepository {
	return &TenantPSQLRepository{}
}

func (repo *TenantPSQLRepository) BeginTx(option ...*sql.TxOptions) (*sql.Tx, error) {
	var op *sql.TxOptions
	if len(option) > 1 {
		op = option[0]
	} else {
		op = &sql.TxOptions{
			Isolation: 0,
			ReadOnly:  false,
		}
	}
	return boil.BeginTx(context.TODO(), op)
}

func (repo *TenantPSQLRepository) InsertTx(tx *sql.Tx, tenant *domain.Tenant) (*domain.Tenant, error) {
	ormTenant := domainTenantToORM(tenant)

	if err := ormTenant.Insert(tx, boil.Infer()); err != nil {
		return nil, err
	}

	return ormTenantToDomain(ormTenant), nil
}

func (repo *TenantPSQLRepository) Update(tenant *domain.Tenant) error {
	ormTenant := domainTenantToORM(tenant)

	rows, err := ormTenant.UpdateG(boil.Infer())

	if err != nil {
		return err
	}
	if rows == 0 {
		return codes.ErrTenantNotFound
	}

	return nil
}

func (repo *TenantPSQLRepository) Delete(id int64) error {
	ormTenant := orm.Tenant{
		ID: id,
	}
	rows, err := ormTenant.DeleteG()

	if err != nil {
		return err
	}
	if rows == 0 {
		return codes.ErrTenantNotFound
	}
	return nil
}

func (repo *TenantPSQLRepository) List(query *domain.TenantQuery) ([]*domain.Tenant, error) {
	mods := make([]qm.QueryMod, 0, 4)

	mods = append(mods, orm.TenantWhere.CreatorID.EQ(query.CreatorID))

	if query.Keyword != "" {
		like := "%" + query.Keyword + "%"
		mods = append(mods, qm.Where(fmt.Sprintf("(%s LIKE ? OR %s LIKE ?)", orm.TenantColumns.Name, orm.TenantColumns.Description), like, like))
	}

	// 游标
	if query.LastID > 0 {
		mods = append(mods, orm.TenantWhere.ID.GT(query.LastID))
	}

	// limit
	mods = append(mods, qm.Limit(query.PageSize))

	// 3.查询数据
	tenants, err := orm.Tenants(mods...).AllG()
	if err != nil {
		return nil, err
	}

	return ormTenantsToDomain(tenants), nil
}

func (repo *TenantPSQLRepository) ExistSameName(creatorID int64, name string) (bool, error) {
	exist, err := orm.Tenants(
		orm.TenantWhere.CreatorID.EQ(creatorID),
		orm.TenantWhere.Name.EQ(name),
	).ExistsG()

	if err != nil {
		return false, errors.WithStack(err)
	}

	return exist, nil
}
