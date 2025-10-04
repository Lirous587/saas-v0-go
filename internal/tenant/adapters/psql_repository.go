package adapters

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/pkg/errors"
	"saas/internal/common/orm"
	"saas/internal/common/reskit/codes"
	"saas/internal/common/utils"
	"saas/internal/tenant/domain"
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

func (repo *TenantPSQLRepository) FindByID(id int64) (*domain.Tenant, error) {
	ormTenant, err := orm.FindTenantG(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, codes.ErrTenantNotFound
		}
		return nil, err
	}
	return ormTenantToDomain(ormTenant), nil
}

func (repo *TenantPSQLRepository) FindTenantPlanByID(id int64) (*domain.Plan, error) {
	// 1.从tenant_plan中查询到plan_id
	tp, err := orm.TenantPlans(
		qm.Where(fmt.Sprintf("%s = ?", orm.TenantPlanColumns.TenantID), id),
		qm.Load(orm.TenantPlanRels.Plan),
	).OneG()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, codes.ErrPlanNotFound
		}
		return nil, err
	}

	if tp.R == nil || tp.R.Plan == nil {
		return nil, codes.ErrPlanNotFound
	}

	return ormPlanToDomain(tp.R.Plan), nil
}

func (repo *TenantPSQLRepository) InsertTx(tx *sql.Tx, tenant *domain.Tenant) (*domain.Tenant, error) {
	ormTenant := domainTenantToORM(tenant)

	if err := ormTenant.Insert(tx, boil.Infer()); err != nil {
		return nil, err
	}

	return ormTenantToDomain(ormTenant), nil
}

func (repo *TenantPSQLRepository) Update(tenant *domain.Tenant) (*domain.Tenant, error) {
	ormTenant := domainTenantToORM(tenant)

	rows, err := ormTenant.UpdateG(boil.Infer())

	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, codes.ErrTenantNotFound
	}

	return ormTenantToDomain(ormTenant), nil
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

func (repo *TenantPSQLRepository) List(query *domain.TenantQuery) (*domain.TenantList, error) {
	var whereMods []qm.QueryMod
	if query.Keyword != "" {
		like := "%" + query.Keyword + "%"
		whereMods = append(whereMods, qm.Where(fmt.Sprintf("(%s LIKE ? OR %s LIKE ?)", orm.TenantColumns.Name, orm.TenantColumns.Description), like, like))
	}
	// 1.计算total
	total, err := orm.Tenants(whereMods...).CountG()
	if err != nil {
		return nil, err
	}

	// 2.计算offset
	offset, err := utils.ComputeOffset(query.Page, query.PageSize)
	if err != nil {
		return nil, err
	}

	listMods := append(whereMods, qm.Offset(offset), qm.Limit(query.PageSize))

	// 3.查询数据
	tenant, err := orm.Tenants(listMods...).AllG()
	if err != nil {
		return nil, err
	}

	return &domain.TenantList{
		Total: total,
		List:  ormTenantsToDomain(tenant),
	}, nil
}

func (repo *TenantPSQLRepository) AssignTenantUserRoleTx(tx *sql.Tx, tenantID, userID, roleID int64) error {
	ormUserTenant := orm.TenantUserRole{
		TenantID: tenantID,
		UserID:   userID,
		RoleID:   roleID,
	}

	return ormUserTenant.Insert(tx, boil.Infer())
}

func (repo *TenantPSQLRepository) AssignTenantUserRole(tenantID, userID, roleID int64) error {
	ormUserTenant := orm.TenantUserRole{
		TenantID: tenantID,
		UserID:   userID,
		RoleID:   roleID,
	}

	return ormUserTenant.InsertG(boil.Infer())
}
