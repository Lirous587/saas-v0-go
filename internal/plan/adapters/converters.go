package adapters

import (
	"saas/internal/common/orm"
	"saas/internal/plan/domain"
	"github.com/aarondl/null/v8"
)

func domainPlanToORM(plan *domain.Plan) *orm.Plan {
	if plan == nil {
		return nil
	}

	// 非null项
	ormPlan := &orm.Plan{
		ID:        		plan.ID,
		Title:     		plan.Title,
        CreatedAt: 		plan.CreatedAt,
        UpdatedAt: 		plan.UpdatedAt,
	}

	// 处理null项
	if plan.Description != "" {
	 	ormPlan.Description = null.StringFrom(plan.Description)
	}

	return ormPlan
}

func ormPlanToDomain(ormPlan *orm.Plan) *domain.Plan {
	if ormPlan == nil {
		return nil
	}

	// 非null项
	plan := &domain.Plan{
		ID:        		ormPlan.ID,
		Title:     		ormPlan.Title,
		CreatedAt: 		ormPlan.CreatedAt,
		UpdatedAt: 		ormPlan.UpdatedAt,
	}

	// 处理null项
	if ormPlan.Description.Valid {
 	 	plan.Description = ormPlan.Description.String
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

