package handler

import (
	"saas/internal/plan/domain"
)

func domainPlanToResponse(plan *domain.Plan) *PlanResponse {
	if plan == nil {
		return nil
	}

	return &PlanResponse{
		ID:          plan.ID,
		Name:        plan.Name,
		Price:       plan.Price,
		Description: plan.Description,
		CreatedAt:   plan.CreatedAt.Unix(),
		UpdatedAt:   plan.UpdatedAt.Unix(),
	}
}

func domainPlansToResponse(plans []*domain.Plan) []*PlanResponse {
	if len(plans) == 0 {
		return nil
	}

	ret := make([]*PlanResponse, 0, len(plans))

	for _, plan := range plans {
		if plan != nil {
			ret = append(ret, domainPlanToResponse(plan))
		}
	}
	return ret
}
