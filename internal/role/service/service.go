package service

import (
	"errors"
	"saas/internal/role/domain"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type service struct {
	repo  domain.RoleRepository
	cache domain.RoleCache
}

func NewRoleService(repo domain.RoleRepository, cache domain.RoleCache) domain.RoleService {
	return &service{
		repo:  repo,
		cache: cache,
	}
}

func (s *service) NewRole() *domain.Role {
	return new(domain.Role)
}

func (s *service) Create(role *domain.Role) (*domain.Role, error) {
	return s.repo.Create(role)
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

func (s *service) GetUserRoleInTenant(userID, tenantID int64) (*domain.Role, error) {
	role, err := s.cache.GetUserRoleInTenant(userID, tenantID)

	if err != nil {
		if errors.Is(err, redis.Nil) {
			role, err = s.repo.FindUserRoleInTenant(userID, tenantID)
			if err != nil {
				return nil, err
			}

			// 添加到缓存
			if err := s.cache.SetUserRoleInTenant(userID, tenantID, role); err != nil {
				zap.L().Error("设置用户角色缓存失败", zap.Error(err), zap.Int64("userID", userID), zap.Int64("tenantID", tenantID))
			}

		} else {
			return nil, err
		}
	}

	return role, nil
}
