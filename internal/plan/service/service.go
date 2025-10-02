package service

import (
	"saas/internal/plan/domain"
)

type service struct {
	repo     domain.PlanRepository
}

func NewPlanService(repo domain.PlanRepository) domain.PlanService {
	return &service{
		repo:     repo,
	}
}

func (s *service) Create(plan *domain.Plan) (*domain.Plan, error) {
	return s.repo.Create(plan)
}

func (s *service) Read(id int64) (*domain.Plan, error) {
   return s.repo.FindByID(id)
}

func (s *service) Update(plan *domain.Plan) (*domain.Plan, error) {
	if _, err := s.repo.FindByID(plan.ID); err != nil {
		return nil, err
	}
	return s.repo.Update(plan)
}

func (s *service) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *service) List(query *domain.PlanQuery) (*domain.PlanList, error) {
	return s.repo.List(query)
}
