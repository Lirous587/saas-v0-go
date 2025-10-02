package service

import (
	"saas/internal/role/domain"
)

type service struct {
	repo     domain.RoleRepository
}

func NewRoleService(repo domain.RoleRepository) domain.RoleService {
	return &service{
		repo:     repo,
	}
}

func (s *service) Create(role *domain.Role) (*domain.Role, error) {
	return s.repo.Create(role)
}

func (s *service) Read(id int64) (*domain.Role, error) {
   return s.repo.FindByID(id)
}

func (s *service) Update(role *domain.Role) (*domain.Role, error) {
	if _, err := s.repo.FindByID(role.ID); err != nil {
		return nil, err
	}
	return s.repo.Update(role)
}

func (s *service) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *service) List(query *domain.RoleQuery) (*domain.RoleList, error) {
	return s.repo.List(query)
}
