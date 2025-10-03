package handler

import (
	"saas/internal/plan/domain"
)

type PlanResponse struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description,omitempty"`
	CreatedAt   int64   `json:"created_at"`
	UpdatedAt   int64   `json:"updated_at"`
}

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

type CreateRequest struct {
	Name        string  `json:"name" binding:"required,max=50"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
}

type UpdateRequest struct {
	Name        string  `json:"name" binding:"required,max=50"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
}

type PlanListResponse struct {
	List []*PlanResponse `json:"list"`
}

func domainPlanListToResponse(data *domain.PlanList) *PlanListResponse {
	if data == nil {
		return nil
	}

	return &PlanListResponse{
		List: domainPlansToResponse(data.List),
	}
}
