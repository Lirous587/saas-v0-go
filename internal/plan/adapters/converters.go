package adapters

import (
	"saas/internal/common/orm"
	"saas/internal/plan/domain"

	"github.com/shopspring/decimal"
)

func domainPlanToORM(plan *domain.Plan) *orm.Plan {
	if plan == nil {
		return nil
	}

	// 非null项
	ormPlan := &orm.Plan{
		ID:          plan.ID,
		Name:        plan.Name,
		Price:       decimal.NewFromFloat(plan.Price),
		Description: plan.Description,
		CreatedAt:   plan.CreatedAt,
		UpdatedAt:   plan.UpdatedAt,
	}

	return ormPlan
}

func ormPlanToDomain(ormPlan *orm.Plan) *domain.Plan {
	if ormPlan == nil {
		return nil
	}

	// 非null项
	plan := &domain.Plan{
		ID:          ormPlan.ID,
		Name:        ormPlan.Name,
		Price:       ormPlan.Price.InexactFloat64(),
		Description: ormPlan.Description,
		CreatedAt:   ormPlan.CreatedAt,
		UpdatedAt:   ormPlan.UpdatedAt,
	}

	return plan
}

func ormPlansToDomain(ormPlans []*orm.Plan) []*domain.Plan {
	if len(ormPlans) == 0 {
		return nil
	}

	plans := make([]*domain.Plan, 0, len(ormPlans))
	for _, ormPlan := range ormPlans {
		if ormPlan != nil {
			plans = append(plans, ormPlanToDomain(ormPlan))
		}
	}
	return plans
}
