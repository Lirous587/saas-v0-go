package service

import (
	"saas/internal/tenant/domain"
)

type service struct {
	repo domain.TenantRepository
}

func NewTenantService(repo domain.TenantRepository) domain.TenantService {
	return &service{
		repo: repo,
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
		return nil, err
	}

	// 2.向tenant_plan插入数据
	if err = s.repo.InsertPlanTx(tx, res.ID, planID); err != nil {
		return nil, err
	}

	// 3.向user_tenant插入数据
	if err = s.repo.InsertUserTx(tx, res.ID, userID); err != nil {
		return nil, err
	}

	// 4.创建此租户的管理员角色

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
