package service

import (
	"saas/internal/comment/domain"
	"saas/internal/common/utils"
)

type service struct {
	repo  domain.CommentRepository
	cache domain.CommentCache
}

func NewCommentService(repo domain.CommentRepository, cache domain.CommentCache) domain.CommentService {
	return &service{
		repo:  repo,
		cache: cache,
	}
}

func (s *service) Create(comment *domain.Comment) (*domain.Comment, error) {
	return s.repo.Create(comment)
}

func (s *service) Read(id int64) (*domain.Comment, error) {
	return s.repo.FindByID(id)
}

func (s *service) Update(comment *domain.Comment) (*domain.Comment, error) {
	if _, err := s.repo.FindByID(comment.ID); err != nil {
		return nil, err
	}
	return s.repo.Update(comment)
}

func (s *service) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *service) List(query *domain.CommentQuery) (*domain.CommentList, error) {
	return s.repo.List(query)
}

func (s *service) SetCommentTenantConfig(config *domain.CommentTenantConfig) error {
	// 生成client_token
	clientToken, err := utils.GenRandomHexToken()
	if err != nil {
		return err
	}

	config.ClientToken = clientToken

	return s.repo.SetCommentTenantConfig(config)
}

func (s *service) GetCommentTenantConfig(tenantID domain.TenantID) (*domain.CommentTenantConfig, error) {
	return s.repo.GetCommentTenantConfig(tenantID)
}

func (s *service) SetCommentConfig(config *domain.CommentConfig) error {
	// 生成client_token
	clientToken, err := utils.GenRandomHexToken()
	if err != nil {
		return err
	}

	config.ClientToken = clientToken

	return s.repo.SetCommentConfig(config)
}

func (s *service) GetCommentConfig(tenantID domain.TenantID, benlongKey domain.BelongKey) (*domain.CommentConfig, error) {
	return s.repo.GetCommentConfig(tenantID, benlongKey)
}
