package service

import (
	"saas/internal/common/reskit/codes"
	"saas/internal/common/utils"
	planDomain "saas/internal/plan/domain"

	"saas/internal/tenant/domain"

	"github.com/friendsofgo/errors"
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

	// 1.检查用户租户限制 (Free/Caring)
	if planID == planDomain.FreePlanID || planID == planDomain.CaringPlanID {
		exist, err := s.planService.CreatorHasPlan(tenant.CreatorID, planID)
		if err != nil {
			return errors.WithMessage(err, "检查用户已有计划失败")
		}

		if exist {
			return codes.ErrPlanUserLimit
		}
	}

	// 2.向tenants插入数据
	res, err := s.repo.InsertTx(tx, tenant)
	if err != nil {
		return errors.WithMessage(err, "向tenants插入数据失败")
	}

	tenantID := res.ID

	// 3.向tenant_plan插入数据
	if err = s.planService.AttchToTenantTx(tx, planID, int64(tenantID), tenant.CreatorID); err != nil {
		return errors.WithMessage(err, "向tenant_plan插入数据失败")
	}

	return tx.Commit()
}

func (s *service) Update(tenant *domain.Tenant) error {
	return s.repo.Update(tenant)
}

func (s *service) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *service) Paging(query *domain.TenantPagingQuery) (*domain.TenantPagination, error) {
	return s.repo.Paging(query)
}

func (s *service) CheckName(creatorID int64, tenantName string) (bool, error) {
	return s.repo.ExistSameName(creatorID, tenantName)
}
