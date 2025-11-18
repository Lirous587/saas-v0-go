package adapters

import (
	"context"
	"database/sql"
	"fmt"
	"saas/internal/common/orm"
	"saas/internal/common/reskit/codes"
	"saas/internal/common/utils/dbkit"
	"saas/internal/tenant/domain"
	"time"

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

func (repo *TenantPSQLRepository) GetByID(id string) (*domain.Tenant, error) {
	ormTenant, err := orm.FindTenantG(
		id,
		orm.TenantColumns.ID,
		orm.TenantColumns.PlanType,
		orm.TenantColumns.Name,
		orm.TenantColumns.Description,
		orm.TenantColumns.CreatedAt,
		orm.TenantColumns.UpdatedAt,
		orm.TenantColumns.CreatorID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, codes.ErrTenantNotFound
		}
	}

	return ormTenantToDomain(ormTenant), nil

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

	rows, err := ormTenant.UpdateG(
		boil.Whitelist(
			orm.TenantColumns.Name,
			orm.TenantColumns.Description,
			orm.TenantColumns.UpdatedAt,
		),
	)

	if err != nil {
		return err
	}
	if rows == 0 {
		return codes.ErrTenantNotFound
	}

	return nil
}

func (repo *TenantPSQLRepository) Delete(id string) error {
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

func (repo *TenantPSQLRepository) ListByKeyset(query *domain.TenantKeysetQuery) (*domain.TenantKeysetResult, error) {
	baseMods := make([]qm.QueryMod, 0, 7)

	// 基本条件
	baseMods = append(baseMods, orm.TenantWhere.CreatorID.EQ(query.CreatorID))

	// 选择列
	baseMods = append(baseMods,
		qm.Select(
			orm.TenantColumns.ID,
			orm.TenantColumns.PlanType,
			orm.TenantColumns.Name,
			orm.TenantColumns.Description,
			orm.TenantColumns.CreatedAt,
			orm.TenantColumns.UpdatedAt,
			orm.TenantColumns.CreatorID,
		),
	)

	if query.Keyword != "" {
		like := "%" + query.Keyword + "%"
		baseMods = append(baseMods, qm.Where(fmt.Sprintf("(%s LIKE ? OR %s LIKE ?)", orm.TenantColumns.Name, orm.TenantColumns.Description), like, like))
	}

	ks := dbkit.NewKeyset[*domain.Tenant](
		orm.TenantColumns.ID,
		orm.TenantColumns.UpdatedAt,
		query.PrevCursor,
		query.NextCursor,
		query.PageSize,
	)

	// 使用 keyset 生成包含 order/limit 的 query mods
	mods := ks.ApplyKeysetMods(baseMods)

	ormTenants, err := orm.Tenants(mods...).AllG()
	if err != nil {
		return nil, err
	}

	domains := ormTenantsToDomain(ormTenants)

	// 精确判断 hasPrev/hasNext：exists 必须和 baseMods 保持一致
	exists := func(primary time.Time, id string, checkPrev bool) (bool, error) {
		var cond qm.QueryMod
		if checkPrev {
			cond = ks.BeforeWhere(primary, id)
		} else {
			cond = ks.AfterWhere(primary, id)
		}
		checkMods := append([]qm.QueryMod{}, baseMods...)
		checkMods = append(checkMods, cond, qm.Limit(1))
		return orm.Tenants(checkMods...).ExistsG()
	}

	result, err := ks.BuildWithExistence(domains, exists)
	if err != nil {
		return nil, err
	}

	return &domain.TenantKeysetResult{
		Items:      result.Items,
		PrevCursor: result.PrevCursor,
		NextCursor: result.NextCursor,
		HasPrev:    result.HasPrev,
		HasNext:    result.HasNext,
	}, nil
}

func (repo *TenantPSQLRepository) ExistSameName(creatorID string, name string) (bool, error) {
	exist, err := orm.Tenants(
		orm.TenantWhere.CreatorID.EQ(creatorID),
		orm.TenantWhere.Name.EQ(name),
	).ExistsG()

	if err != nil {
		return false, errors.WithStack(err)
	}

	return exist, nil
}

func (repo *TenantPSQLRepository) IsCreatorHasPlan(creatorID string, planType domain.PlanType) (bool, error) {
	exist, err := orm.Tenants(
		orm.TenantWhere.CreatorID.EQ(creatorID),
		orm.TenantWhere.PlanType.EQ(orm.TenantPlanType(planType)),
	).ExistsG()

	if err != nil {
		return false, errors.WithStack(err)
	}
	return exist, nil

}

func (repo *TenantPSQLRepository) GetPlan(id string) (*domain.Plan, error) {
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
