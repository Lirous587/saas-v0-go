package adapters

import (
	"database/sql"
	"fmt"
	"saas/internal/comment/domain"
	"saas/internal/common/orm"
	"saas/internal/common/reskit/codes"
	"saas/internal/common/utils"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/pkg/errors"
)

type CommentPSQLRepository struct {
}

func NewCommentPSQLRepository() domain.CommentRepository {
	return &CommentPSQLRepository{}
}

func (repo *CommentPSQLRepository) FindByID(id int64) (*domain.Comment, error) {
	ormComment, err := orm.FindCommentG(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, codes.ErrCommentNotFound
		}
		return nil, err
	}
	return ormCommentToDomain(ormComment), nil
}

func (repo *CommentPSQLRepository) Create(comment *domain.Comment) (*domain.Comment, error) {
	ormComment := domainCommentToORM(comment)

	if err := ormComment.InsertG(boil.Infer()); err != nil {
		return nil, err
	}

	return ormCommentToDomain(ormComment), nil
}

func (repo *CommentPSQLRepository) Update(comment *domain.Comment) (*domain.Comment, error) {
	ormComment := domainCommentToORM(comment)

	rows, err := ormComment.UpdateG(boil.Infer())

	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, codes.ErrCommentNotFound
	}

	return ormCommentToDomain(ormComment), nil
}

func (repo *CommentPSQLRepository) Delete(id int64) error {
	ormComment := orm.Comment{
		ID: id,
	}
	rows, err := ormComment.DeleteG()

	if err != nil {
		return err
	}
	if rows == 0 {
		return codes.ErrCommentNotFound
	}
	return nil
}

func (repo *CommentPSQLRepository) List(query *domain.CommentQuery) (*domain.CommentList, error) {
	var whereMods []qm.QueryMod
	// if query.Keyword != "" {
	// 	// like := "%" + query.Keyword + "%"
	// 	// whereMods = append(whereMods, qm.Where(fmt.Sprintf("(%s LIKE ? OR %s LIKE ?)", orm.CommentColumns.Title, orm.CommentColumns.Description), like, like))
	// }
	// 1.计算total
	total, err := orm.Comments(whereMods...).CountG()
	if err != nil {
		return nil, err
	}

	// 2.计算offset
	offset, err := utils.ComputeOffset(query.Page, query.PageSize)
	if err != nil {
		return nil, err
	}

	listMods := append(whereMods, qm.Offset(offset), qm.Limit(query.PageSize))

	// 3.查询数据
	comment, err := orm.Comments(listMods...).AllG()
	if err != nil {
		return nil, err
	}

	return &domain.CommentList{
		Total: total,
		List:  ormCommentsToDomain(comment),
	}, nil
}

func (repo *CommentPSQLRepository) SetCommentTenantConfig(config *domain.CommentTenantConfig) error {
	ormConfig := domainCommentTenantConfigToORM(config)
	if err := ormConfig.UpsertG(
		true,
		[]string{orm.CommentTenantConfigColumns.TenantID},
		boil.Whitelist(
			orm.CommentTenantConfigColumns.IfAudit,
			orm.CommentTenantConfigColumns.ClientToken,
			orm.CommentTenantConfigColumns.UpdatedAt,
		),
		boil.Infer(),
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (repo *CommentPSQLRepository) GetCommentTenantConfig(tenantID domain.TenantID) (*domain.CommentTenantConfig, error) {
	ormConfig, err := orm.CommentTenantConfigs(
		qm.Where(fmt.Sprintf("%s = ?", orm.CommentTenantConfigColumns.TenantID), tenantID),
	).OneG()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, codes.ErrCommentTenantConfigNotFound
		}
		return nil, err
	}

	return ormCommentTenantConfigToDomain(ormConfig), nil
}

func (repo *CommentPSQLRepository) SetCommentConfig(config *domain.CommentConfig) error {
	ormConfig := domainCommentConfigToORM(config)
	if err := ormConfig.UpsertG(
		true,
		[]string{orm.CommentConfigColumns.TenantID, orm.CommentConfigColumns.BelongKey}, // 冲突列：复合主键的两个字段
		boil.Whitelist(
			orm.CommentConfigColumns.IfAudit,
			orm.CommentConfigColumns.ClientToken,
			orm.CommentConfigColumns.UpdatedAt,
		),
		boil.Infer(),
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (repo *CommentPSQLRepository) GetCommentConfig(tenantID domain.TenantID, benlongKey domain.BelongKey) (*domain.CommentConfig, error) {
	ormConfig, err := orm.CommentConfigs(
		qm.Where(fmt.Sprintf("%s = ?", orm.CommentConfigColumns.TenantID), tenantID),
		qm.Where(fmt.Sprintf("%s = ?", orm.CommentConfigColumns.BelongKey), benlongKey),
	).OneG()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, codes.ErrCommentConfigNotFound
		}
		return nil, err
	}

	return ormCommentConfigToDomain(ormConfig), nil
}
