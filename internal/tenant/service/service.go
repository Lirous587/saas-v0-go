package service

import (
	planDomain "saas/internal/plan/domain"
	roleDomain "saas/internal/role/domain"
	"saas/internal/tenant/domain"

	"github.com/friendsofgo/errors"
)

type service struct {
	repo        domain.TenantRepository
	planService planDomain.PlanService
	roleService roleDomain.RoleService
}

func NewTenantService(repo domain.TenantRepository, planService planDomain.PlanService, roleService roleDomain.RoleService) domain.TenantService {
	return &service{
		repo:        repo,
		planService: planService,
		roleService: roleService,
	}
}

func (s *service) Create(tenant *domain.Tenant, planID int64, userID int64) (*domain.Tenant, error) {
	// 1.创建事务
	tx, err := s.repo.BeginTx()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 1.向tenants插入数据
	res, err := s.repo.InsertTx(tx, tenant)
	if err != nil {
		return nil, errors.WithMessage(err, "向tenants插入数据失败")
	}

	tenantID := res.ID

	// 2.向tenant_plan插入数据
	if err = s.planService.AttchToTenantTx(tx, planID, tenantID); err != nil {
		return nil, errors.WithMessage(err, "向tenant_plan插入数据失败")
	}

	// 3.为此租户设置superadmin
	superadmin := s.roleService.NewRole().GetDefultSuperadmin()

	if err = s.repo.AssignTenantUserRoleTx(tx, tenantID, userID, superadmin.ID); err != nil {
		return nil, errors.WithMessage(err, "为此租户设置superadmin失败")
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return res, nil
}

func (s *service) Read(id int64) (*domain.Tenant, error) {
	return s.repo.FindByID(id)
}

func (s *service) Update(tenant *domain.Tenant) (*domain.Tenant, error) {
	if _, err := s.repo.FindByID(tenant.ID); err != nil {
		return nil, err
	}
	return s.repo.Update(tenant)
}

func (s *service) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *service) List(query *domain.TenantQuery) (*domain.TenantList, error) {
	return s.repo.List(query)
}
