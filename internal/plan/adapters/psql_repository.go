package adapters

import (
	"database/sql"
    "fmt"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/pkg/errors"
	"saas/internal/common/reskit/codes"
	"saas/internal/plan/domain"
	"saas/internal/common/orm"
	"saas/internal/common/utils"
)

type PSQLPlanRepository struct {
}

func NewPSQLPlanRepository() domain.PlanRepository {
	return &PSQLPlanRepository{}
}

func (repo *PSQLPlanRepository) FindByID(id int64) (*domain.Plan, error) {
	ormPlan, err := orm.FindPlanG(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, codes.ErrPlanNotFound
		}
		return nil, err
	}
	return ormPlanToDomain(ormPlan), nil
}

func (repo *PSQLPlanRepository) Create(plan *domain.Plan) (*domain.Plan,error)  {
	ormPlan := domainPlanToORM(plan)

	if err := ormPlan.InsertG(boil.Infer()); err != nil {
		return nil, err
	}

	return ormPlanToDomain(ormPlan), nil
}

func (repo *PSQLPlanRepository) Update(plan *domain.Plan) (*domain.Plan,error) {
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

func (repo *PSQLPlanRepository) Delete(id int64) error {
	ormPlan := orm.Plan{
		ID: id,
	}
	rows, err := ormPlan.DeleteG(false)

	if err != nil {
		return err
	}
	if rows == 0 {
		return codes.ErrPlanNotFound
	}
	return nil
}

func (repo *PSQLPlanRepository) List(query *domain.PlanQuery) (*domain.PlanList, error) {
	var whereMods []qm.QueryMod
	if query.Keyword != "" {
		like := "%" + query.Keyword + "%"
		whereMods = append(whereMods, qm.Where(fmt.Sprintf("(%s LIKE ? OR %s LIKE ?)", orm.PlanColumns.Title, orm.PlanColumns.Description), like, like))
	}
	// 1.计算total
	total, err := orm.Plans(whereMods...).CountG()
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
	plan, err := orm.Plans(listMods...).AllG()
	if err != nil {
		return nil, err
	}

	return &domain.PlanList{
		Total: total,
		List:  ormPlansToDomain(plan),
	}, nil
}
