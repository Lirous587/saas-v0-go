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

func (repo *CommentPSQLRepository) CreatePlate(plate *domain.Plate) error {
	ormPlate := domainPlateToORM(plate)

	return ormPlate.InsertG(boil.Infer())
}

func (repo *CommentPSQLRepository) DeletePlate(tenantID domain.TenantID, id int64) error {
	rows, err := orm.CommentPlates(
		qm.Where(fmt.Sprintf("%s = ?", orm.CommentPlateColumns.TenantID), tenantID),
		qm.Where(fmt.Sprintf("%s = ?", orm.CommentPlateColumns.ID), id),
	).DeleteAllG()

	if err != nil {
		return err
	}

	if rows == 0 {
		return codes.ErrCommentPlateNotFound
	}

	return nil
}

func (repo *CommentPSQLRepository) ListPlate(query *domain.PlateQuery) (*domain.PlateList, error) {
	var whereMods []qm.QueryMod
	if query.Keyword != "" {
		like := "%" + query.Keyword + "%"
		whereMods = append(whereMods, qm.Where(fmt.Sprintf("(%s LIKE ? OR %s LIKE ?)", orm.CommentPlateColumns.BelongKey, orm.CommentPlateColumns.Summary), like, like))
	}
	// 1.计算total
	total, err := orm.CommentPlates(whereMods...).CountG()
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
	plate, err := orm.CommentPlates(listMods...).AllG()
	if err != nil {
		return nil, err
	}

	return &domain.PlateList{
		Total: total,
		List:  ormPlatesToDomain(plate),
	}, nil
}

func (repo *CommentPSQLRepository) ExistPlateBykey(tenantID domain.TenantID, belongKey string) (bool, error) {
	exist, err := orm.CommentPlates(
		qm.Where(fmt.Sprintf("%s = ?", orm.CommentPlateColumns.TenantID), tenantID),
		qm.Where(fmt.Sprintf("%s = ?", orm.CommentPlateColumns.BelongKey), belongKey),
	).ExistsG()

	if err != nil {
		return false, err
	}

	return exist, nil
}
func (repo *CommentPSQLRepository) GetPlateBelongByID(id int64) (*domain.PlateBelong, error) {
	plate, err := orm.CommentPlates(
		qm.Where(fmt.Sprintf("%s = ?", orm.CommentPlateColumns.ID), id),
		qm.Select(orm.CommentPlateColumns.ID, orm.CommentPlateColumns.BelongKey),
	).OneG()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, codes.ErrCommentPlateNotFound
		}
		return nil, err
	}

	return &domain.PlateBelong{
		ID:        plate.ID,
		BelongKey: plate.BelongKey,
	}, nil
}

func (repo *CommentPSQLRepository) GetPlateBelongByKey(tenantID domain.TenantID, belongKey string) (*domain.PlateBelong, error) {
	plate, err := orm.CommentPlates(
		qm.Where(fmt.Sprintf("%s = ?", orm.CommentPlateColumns.TenantID), tenantID),
		qm.Where(fmt.Sprintf("%s = ?", orm.CommentPlateColumns.BelongKey), belongKey),
		qm.Select(orm.CommentPlateColumns.ID, orm.CommentPlateColumns.BelongKey),
	).OneG()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, codes.ErrCommentPlateNotFound
		}
		return nil, err
	}

	return &domain.PlateBelong{
		ID:        plate.ID,
		BelongKey: plate.BelongKey,
	}, nil
}

func (repo *CommentPSQLRepository) SetTenantConfig(config *domain.TenantConfig) error {
	ormConfig := domainTenantConfigToORM(config)
	if err := ormConfig.UpsertG(
		true,
		[]string{orm.CommentTenantConfigColumns.TenantID},
		boil.Greylist( // 手动指定 insertColumns，确保包含 if_audit
			orm.CommentTenantConfigColumns.IfAudit,
		),
		boil.Greylist( // 手动指定 insertColumns，确保包含 if_audit
			orm.CommentTenantConfigColumns.IfAudit,
		),
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (repo *CommentPSQLRepository) GetTenantConfig(tenantID domain.TenantID) (*domain.TenantConfig, error) {
	ormConfig, err := orm.CommentTenantConfigs(
		qm.Where(fmt.Sprintf("%s = ?", orm.CommentTenantConfigColumns.TenantID), tenantID),
	).OneG()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, codes.ErrCommentTenantConfigNotFound
		}
		return nil, err
	}

	return ormTenantConfigToDomain(ormConfig), nil
}

func (repo *CommentPSQLRepository) SetPlateConfig(config *domain.PlateConfig) error {
	ormConfig := domainPlateConfigToORM(config)
	if err := ormConfig.UpsertG(
		true,
		[]string{orm.CommentPlateConfigColumns.TenantID, orm.CommentPlateConfigColumns.PlateID}, // 冲突列：复合主键的两个字段
		boil.Greylist(
			orm.CommentPlateConfigColumns.IfAudit,
		),
		boil.Greylist(
			orm.CommentPlateConfigColumns.IfAudit,
		),
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (repo *CommentPSQLRepository) GetPlateConfig(tenantID domain.TenantID, plateID int64) (*domain.PlateConfig, error) {
	ormConfig, err := orm.CommentPlateConfigs(
		qm.Where(fmt.Sprintf("%s = ?", orm.CommentPlateConfigColumns.TenantID), tenantID),
		qm.Where(fmt.Sprintf("%s = ?", orm.CommentPlateConfigColumns.PlateID), plateID),
	).OneG()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, codes.ErrCommentPlateConfigNotFound
		}
		return nil, err
	}

	return ormPlateConfigToDomain(ormConfig), nil
}
