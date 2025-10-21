package service

import (
	"saas/internal/common/utils"
	planDomain "saas/internal/plan/domain"

	"github.com/friendsofgo/errors"
	"saas/internal/tenant/domain"
)

type service struct {
	myDomain    string
	repo        domain.TenantRepository
	cache       domain.TenantCache
	planService planDomain.PlanService
}

func NewTenantService(repo domain.TenantRepository, cache domain.TenantCache, planService planDomain.PlanService) domain.TenantService {
	myDomain := utils.GetEnv("DOMAIN")

	return &service{
		myDomain:    myDomain,
		repo:        repo,
		cache:       cache,
		planService: planService,
	}
}

func (s *service) Create(tenant *domain.Tenant, planID int64) error {
	// 1.创建事务
	tx, err := s.repo.BeginTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1.向tenants插入数据
	res, err := s.repo.InsertTx(tx, tenant)
	if err != nil {
		return errors.WithMessage(err, "向tenants插入数据失败")
	}

	tenantID := res.ID

	// 2.向tenant_plan插入数据
	if err = s.planService.AttchToTenantTx(tx, planID, tenantID); err != nil {
		return errors.WithMessage(err, "向tenant_plan插入数据失败")
	}

	return tx.Commit()
}

func (s *service) Read(id int64) (*domain.Tenant, error) {
	return s.repo.FindByID(id)
}

func (s *service) Update(tenant *domain.Tenant) error {
	return s.repo.Update(tenant)
}

func (s *service) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *service) List(query *domain.TenantQuery) (*domain.TenantList, error) {
	return s.repo.List(query)
}
