package service

import (
	"saas/internal/common/reskit/codes"

	"saas/internal/tenant/domain"

	"github.com/friendsofgo/errors"
)

type service struct {
	repo  domain.TenantRepository
	cache domain.TenantCache
}

func NewTenantService(repo domain.TenantRepository, cache domain.TenantCache) domain.TenantService {
	return &service{
		repo:  repo,
		cache: cache,
	}
}

func (s *service) Create(tenant *domain.Tenant) error {
	// 检查用户的租户计划限制 (Free/Care)
	if tenant.PlanType == domain.PlanFreeType || tenant.PlanType == domain.PlanCareType {
		exist, err := s.repo.IsCreatorHasPlan(tenant.CreatorID, tenant.PlanType)
		if err != nil {
			return errors.WithMessage(err, "检查用户已有计划失败")
		}

		if exist {
			return codes.ErrTenantPlanUserLimit
		}
	}

	// 向tenants插入数据
	_, err := s.repo.Create(tenant)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *service) Update(tenant *domain.Tenant) error {
	return s.repo.Update(tenant)
}

func (s *service) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *service) GetByID(id int64) (*domain.Tenant, error) {
	return s.repo.GetByID(id)
}

func (s *service) Paging(query *domain.TenantPagingQuery) (*domain.TenantPagination, error) {
	return s.repo.Paging(query)
}

func (s *service) CheckName(creatorID int64, tenantName string) (bool, error) {
	return s.repo.ExistSameName(creatorID, tenantName)
}

func (s *service) GetPlan(id int64) (*domain.Plan, error) {
	return s.repo.GetPlan(id)
}
