package service

import (
	"database/sql"
	"saas/internal/plan/domain"
)

type service struct {
	repo domain.PlanRepository
}

func NewPlanService(repo domain.PlanRepository) domain.PlanService {
	return &service{
		repo: repo,
	}
}

func (s *service) Create(plan *domain.Plan) error {
	return s.repo.Create(plan)
}

func (s *service) Update(plan *domain.Plan) error {
	return s.repo.Update(plan)
}

func (s *service) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *service) List() ([]*domain.Plan, error) {
	return s.repo.List()
}

func (s *service) AttchToTenantTx(tx *sql.Tx, planID, tenantID int64) error {
	return s.repo.AttchToTenantTx(tx, planID, tenantID)
}
