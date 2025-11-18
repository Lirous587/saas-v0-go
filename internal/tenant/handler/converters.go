package handler

import (
	"saas/internal/tenant/domain"
)

func domainTenantToResponse(tenant *domain.Tenant) *TenantResponse {
	if tenant == nil {
		return nil
	}

	return &TenantResponse{
		ID:          tenant.ID,
		PlanType:    tenant.PlanType,
		CreatorID:   tenant.CreatorID,
		Name:        tenant.Name,
		Description: tenant.Description,
		CreatedAt:   tenant.CreatedAt.Unix(),
		UpdatedAt:   tenant.UpdatedAt.Unix(),
	}
}

func domainTenantsToResponse(tenants []*domain.Tenant) []*TenantResponse {
	if len(tenants) == 0 {
		return nil
	}

	ret := make([]*TenantResponse, 0, len(tenants))

	for _, tenant := range tenants {
		if tenant != nil {
			ret = append(ret, domainTenantToResponse(tenant))
		}
	}
	return ret
}

func domainTenantKeysetToResponse(pager *domain.TenantKeysetResult) *ListByKeysetResponse {
	if pager == nil {
		return nil
	}

	return &ListByKeysetResponse{
		Items:      domainTenantsToResponse(pager.Items),
		PrevCursor: pager.PrevCursor,
		NextCursor: pager.NextCursor,
		HasPrev:    pager.HasPrev,
		HasNext:    pager.HasNext,
	}
}

func domainPlanToResponse(plan *domain.Plan) *PlanResponse {
	if plan == nil {
		return nil
	}

	return &PlanResponse{
		TenantID:     plan.TenantID,
		PlanType:     plan.PlanType,
		StartTime:    plan.StartTime.Unix(),
		EndTime:      plan.EndTime.Unix(),
		Status:       plan.Status,
		BillingCycle: plan.BillingCycle,
		CanUpgrade:   plan.CanUpgrade,
	}
}
