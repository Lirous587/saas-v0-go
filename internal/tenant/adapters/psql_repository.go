package adapters

import (
	"context"
	"database/sql"
	"fmt"
	"saas/internal/common/orm"
	"saas/internal/common/reskit/codes"
	"saas/internal/common/utils/dbkit"
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

func (repo *TenantPSQLRepository) Paging(query *domain.TenantPagingQuery) (*domain.TenantPagination, error) {
	mods := make([]qm.QueryMod, 0, 5)

	// 基本条件
	mods = append(mods, orm.TenantWhere.CreatorID.EQ(query.CreatorID))

	if query.Keyword != "" {
		like := "%" + query.Keyword + "%"
		mods = append(mods, qm.Where(fmt.Sprintf("(%s LIKE ? OR %s LIKE ?)", orm.TenantColumns.Name, orm.TenantColumns.Description), like, like))
	}

	keyset := dbkit.NewKeyset[domain.Tenant](query.PageSize, query.BeforeID, query.AfterID)

	mods = keyset.ApplyKeysetMods(mods, orm.TenantColumns.ID)

	ormTenants, err := orm.Tenants(mods...).AllG()
	if err != nil {
		return nil, err
	}

	domains := ormTenantsToDomain(ormTenants)

	result := keyset.BuildPaginationResult(domains)

	return &domain.TenantPagination{
		Items:   result.Items,
		HasNext: result.HasNext,
		HasPrev: result.HasPrev,
	}, nil
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
