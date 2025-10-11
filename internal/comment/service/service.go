package service

import (
	"saas/internal/comment/domain"
	"saas/internal/common/email"
	"saas/internal/common/reskit/codes"
	"saas/internal/common/utils"

	"github.com/pkg/errors"
)

type service struct {
	repo   domain.CommentRepository
	cache  domain.CommentCache
	mailer email.Mailer
}

func NewCommentService(repo domain.CommentRepository, cache domain.CommentCache, mailer email.Mailer) domain.CommentService {
	return &service{
		repo:   repo,
		cache:  cache,
		mailer: mailer,
	}
}

func (s *service) Create(comment *domain.Comment) (*domain.Comment, error) {
	// 1.plate 是否存在
	exist, err := s.repo.ExistPlateBykey(comment.TenantID, comment.Plate.BelongKey)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if !exist {
		return nil, errors.WithStack(codes.ErrCommentPlateNotFound)
	}

	// 2.检查parent_id和root_id 根据其来发送邮件
	if comment.HasPartent() {

	}

	// user_id
	s.repo.Create(comment)

	return nil, nil
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

func (s *service) CreatePlate(plate *domain.Plate) error {
	exist, err := s.repo.ExistPlateBykey(plate.TenantID, plate.BelongKey)
	if err != nil {
		return errors.WithStack(err)
	}

	if exist {
		return codes.ErrCommentPlateExist
	}

	if err := s.repo.CreatePlate(plate); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s *service) DeletePlate(tenantID domain.TenantID, id int64) error {
	return s.repo.DeletePlate(tenantID, id)
}

func (s *service) ListPlate(query *domain.PlateQuery) (*domain.PlateList, error) {
	return s.repo.ListPlate(query)
}

func (s *service) SetTenantConfig(config *domain.TenantConfig) error {
	// 生成client_token
	clientToken, err := utils.GenRandomHexToken()
	if err != nil {
		return err
	}

	config.ClientToken = clientToken

	return s.repo.SetTenantConfig(config)
}

func (s *service) GetTenantConfig(tenantID domain.TenantID) (*domain.TenantConfig, error) {
	return s.repo.GetTenantConfig(tenantID)
}

func (s *service) SetPlateConfig(config *domain.PlateConfig) error {
	plate, err := s.repo.GetPlateBelongByKey(config.TenantID, config.Plate.BelongKey)
	if err != nil {
		return errors.WithStack(err)
	}

	config.Plate.ID = plate.ID
	if err := s.repo.SetPlateConfig(config); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *service) GetPlateConfig(tenantID domain.TenantID, belongKey string) (*domain.PlateConfig, error) {
	plate, err := s.repo.GetPlateBelongByKey(tenantID, belongKey)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	config, err := s.repo.GetPlateConfig(tenantID, plate.ID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return config, nil
}
