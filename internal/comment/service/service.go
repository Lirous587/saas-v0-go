package service

import (
	"saas/internal/comment/domain"
)

type service struct {
	repo     domain.CommentRepository
}

func NewCommentService(repo domain.CommentRepository) domain.CommentService {
	return &service{
		repo:     repo,
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
