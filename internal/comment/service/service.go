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

func (s *service) Audit(tenantID domain.TenantID, id int64, status domain.CommentStatus) error {
	comment, err := s.repo.GetByID(tenantID, id)
	if err != nil {
		return errors.WithStack(err)
	}

	if !comment.CanAudit() {
		return codes.ErrCommentIllegalAudit
	}

	// 同意
	if comment.IsApproved() {
		if err := s.repo.Approve(tenantID, id); err != nil {
			return errors.WithMessage(err, "同意评论时候更新status失败")
		}
	} else {
		if err := s.repo.Delete(tenantID, id); err != nil {
			return errors.WithMessage(err, "拒绝评论时候删除评论记录失败")
		}
	}

	// 通知评论者

	return nil
}

func (s *service) Delete(tenantID domain.TenantID, userID int64, id int64) error {
	// 查询当前评论用户
	uid, err := s.repo.GetCommentUser(tenantID, id)
	if err != nil {
		return errors.WithStack(err)
	}

	// 如果请求用户和评论用户不一致
	if uid != userID {
		// 去获取当前租户的uid
		admin, err := s.repo.GetDomainAdminByTenant(tenantID)
		if err != nil {
			return errors.WithStack(err)
		}

		if userID != admin.ID {
			return codes.ErrCommentNoPermissionToDelete
		}
	}

	return s.repo.Delete(tenantID, id)
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

func (s *service) GetPlateConfig(tenantID domain.TenantID, plateID int64) (*domain.PlateConfig, error) {
	config, err := s.repo.GetPlateConfig(tenantID, plateID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return config, nil
}
