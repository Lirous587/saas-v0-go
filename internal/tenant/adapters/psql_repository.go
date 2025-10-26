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

func (repo *TenantPSQLRepository) Create(tenant *domain.Tenant) (*domain.Tenant, error) {
	ormTenant := domainTenantToORM(tenant)

	if err := ormTenant.InsertG(boil.Infer()); err != nil {
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

	keyset := dbkit.NewKeyset[domain.Tenant](orm.TenantColumns.ID, orm.TenantColumns.CreatedAt, query.PrevCursor, query.NextCursor, query.PageSize)

	mods = keyset.ApplyKeysetMods(mods)

	ormTenants, err := orm.Tenants(mods...).AllG()
	if err != nil {
		return nil, err
	}

	domains := ormTenantsToDomain(ormTenants)

	result := keyset.BuildPaginationResult(domains)

	return &domain.TenantPagination{
		Items:      result.Items,
		PrevCursor: result.PrevCursor,
		NextCursor: result.NextCursor,
		HasPrev:    result.HasPrev,
		HasNext:    result.HasNext,
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

func (repo *TenantPSQLRepository) IsCreatorHasPlan(creatorID int64, planType domain.PlanType) (bool, error) {
	exist, err := orm.Tenants(
		orm.TenantWhere.CreatorID.EQ(creatorID),
		orm.TenantWhere.PlanType.EQ(orm.TenantPlanType(planType)),
	).ExistsG()

	if err != nil {
		return false, errors.WithStack(err)
	}
	return exist, nil

}

func (repo *TenantPSQLRepository) GetPlan(id int64) (*domain.Plan, error) {
	tenantPlan, err := orm.Tenants(
		orm.TenantWhere.ID.EQ(id),
	).OneG()

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, codes.ErrTenantPlanNotFound
		}
		return nil, errors.WithStack(err)
	}

	return ormTenantPlanToDomain(tenantPlan), nil
}
