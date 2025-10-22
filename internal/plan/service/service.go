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
	if _, err := s.repo.Create(plan); err != nil {
		return err
	}
	return nil
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

func (s *service) CreatorHasPlan(creatorID, planID int64) (bool, error) {
	return s.repo.CreatorHasPlan(creatorID, planID)
}

func (s *service) AttchToTenantTx(tx *sql.Tx, planID, tenantID, creatorID int64) error {
	return s.repo.AttchToTenantTx(tx, planID, tenantID, creatorID)
}
