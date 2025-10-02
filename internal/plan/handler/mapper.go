package handler

import (
    "saas/internal/plan/domain"
)

type PlanResponse struct {
    ID          int64  `json:"id"`
    Title       string `json:"title"`
    Description string `json:"description,omitempty"`
    CreatedAt   int64  `json:"created_at"`
    UpdatedAt   int64  `json:"updated_at"`
}

func domainPlanToResponse(plan *domain.Plan) *PlanResponse {
    if plan == nil {
        return nil
    }

    return &PlanResponse{
        ID:          plan.ID,
        Title:       plan.Title,
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
    Title       string  `json:"title" binding:"required,max=30"`
    Description string  `json:"description" binding:"max=60"`
}

type UpdateRequest struct {
    Title       string  `json:"title" binding:"required,max=30"`
    Description string  `json:"description" binding:"max=60"`
}

type ListRequest struct {
    Page     int    `form:"page,default=1" binding:"min=1"`
    PageSize int    `form:"page_size,default=5" binding:"min=5,max=20"`
    KeyWord  string `form:"keyword" binding:"max=20"`
}

type PlanListResponse struct {
    Total int64                         `json:"total"`
    List  []*PlanResponse   `json:"list"`
}

func domainPlanListToResponse(data *domain.PlanList) *PlanListResponse {
    if data == nil {
        return nil
    }

    return &PlanListResponse{
        Total: data.Total,
        List:  domainPlansToResponse(data.List),
    }
}
