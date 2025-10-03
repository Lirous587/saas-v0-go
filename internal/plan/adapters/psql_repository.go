package adapters

import (
	"database/sql"
	"fmt"
	"saas/internal/common/orm"
	"saas/internal/common/reskit/codes"
	"saas/internal/plan/domain"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/pkg/errors"
)

type PlanPSQLRepository struct {
}

func NewPlanPSQLRepository() domain.PlanRepository {
	return &PlanPSQLRepository{}
}

func (repo *PlanPSQLRepository) FindByID(id int64) (*domain.Plan, error) {
	ormPlan, err := orm.FindPlanG(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, codes.ErrPlanNotFound
		}
		return nil, err
	}
	return ormPlanToDomain(ormPlan), nil
}

func (repo *PlanPSQLRepository) Create(plan *domain.Plan) (*domain.Plan, error) {
	ormPlan := domainPlanToORM(plan)

	if err := ormPlan.InsertG(boil.Infer()); err != nil {
		return nil, err
	}

	return ormPlanToDomain(ormPlan), nil
}

func (repo *PlanPSQLRepository) Update(plan *domain.Plan) (*domain.Plan, error) {
	ormPlan := domainPlanToORM(plan)

	rows, err := ormPlan.UpdateG(boil.Infer())

	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, codes.ErrPlanNotFound
	}

	return ormPlanToDomain(ormPlan), nil
}

func (repo *PlanPSQLRepository) Delete(id int64) error {
	ormPlan := orm.Plan{
		ID: id,
	}
	rows, err := ormPlan.DeleteG()

	if err != nil {
		return err
	}
	if rows == 0 {
		return codes.ErrPlanNotFound
	}
	return nil
}

func (repo *PlanPSQLRepository) List() (*domain.PlanList, error) {
	plan, err := orm.Plans(qm.OrderBy(fmt.Sprintf("%s ASC", orm.PlanColumns.Price))).AllG()
	if err != nil {
		return nil, err
	}

	return &domain.PlanList{
		List: ormPlansToDomain(plan),
	}, nil
}

func (repo *PlanPSQLRepository) AttchToTenantTx(tx *sql.Tx, planID, tenantID int64) error {
	ormTenantPlan := orm.TenantPlan{
		TenantID: tenantID,
		PlanID:   planID,
	}

	return ormTenantPlan.Insert(tx, boil.Infer())
}
