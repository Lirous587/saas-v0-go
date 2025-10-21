package adapters

import (
	"github.com/aarondl/null/v8"
	"saas/internal/common/orm"
	"saas/internal/tenant/domain"
)

func domainTenantToORM(tenant *domain.Tenant) *orm.Tenant {
	if tenant == nil {
		return nil
	}

	// 非null项
	ormTenant := &orm.Tenant{
		ID:        tenant.ID,
		Name:      tenant.Name,
		CreatedAt: tenant.CreatedAt,
		UpdatedAt: tenant.UpdatedAt,
	}

	// 处理null项
	if tenant.Description != "" {
		ormTenant.Description = null.StringFrom(tenant.Description)
	}

	return ormTenant
}

func ormTenantToDomain(ormTenant *orm.Tenant) *domain.Tenant {
	if ormTenant == nil {
		return nil
	}

	// 非null项
	tenant := &domain.Tenant{
		ID:        ormTenant.ID,
		Name:      ormTenant.Name,
		CreatedAt: ormTenant.CreatedAt,
		UpdatedAt: ormTenant.UpdatedAt,
	}

	// 处理null项
	if ormTenant.Description.Valid {
		tenant.Description = ormTenant.Description.String
	}

	return tenant
}

func ormTenantsToDomain(ormTenants []*orm.Tenant) []*domain.Tenant {
	if len(ormTenants) == 0 {
		return nil
	}

	tenants := make([]*domain.Tenant, 0, len(ormTenants))
	for _, ormTenant := range ormTenants {
		if ormTenant != nil {
			tenants = append(tenants, ormTenantToDomain(ormTenant))
		}
	}
	return tenants
}

func ormPlanToDomain(ormPlan *orm.Plan) *domain.Plan {
	if ormPlan == nil {
		return nil
	}

	// 非null项
	plan := &domain.Plan{
		ID:   ormPlan.ID,
		Name: ormPlan.Name,
	}

	return plan
}
