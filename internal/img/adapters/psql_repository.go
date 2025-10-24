package adapters

import (
	"database/sql"
	"fmt"
	"saas/internal/common/orm"
	"saas/internal/common/reskit/codes"
	"saas/internal/common/utils/dbkit"
	"saas/internal/img/domain"
	"time"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/pkg/errors"
)

type ImgPSQLRepository struct {
}

func NewImgPSQLRepository() domain.ImgRepository {
	return &ImgPSQLRepository{}
}

func (repo *ImgPSQLRepository) FindByID(tenantID domain.TenantID, id int64, deleted ...bool) (*domain.Img, error) {
	selectDeleted := len(deleted) > 0 && deleted[0]

	var whereMods []qm.QueryMod
	whereMods = append(whereMods,
		qm.Where(fmt.Sprintf("%s = ?", orm.ImgColumns.TenantID), tenantID),
		qm.And(fmt.Sprintf("%s = ?", orm.ImgColumns.ID), id),
	)

	if selectDeleted {
		whereMods = append(whereMods, qm.WithDeleted())
	}

	ormImg, err := orm.Imgs(whereMods...).OneG()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, codes.ErrImgNotFound
		}
		return nil, err
	}

	return ormImgToDomain(ormImg), err
}

func (repo *ImgPSQLRepository) ExistByPath(tenantID domain.TenantID, path string) (bool, error) {
	exist, err := orm.Imgs(
		qm.Where(fmt.Sprintf("%s = ?", orm.ImgColumns.TenantID), tenantID),
		qm.And(fmt.Sprintf("%s = ?", orm.ImgColumns.Path), path),
		qm.WithDeleted(),
	).ExistsG()
	if err != nil {
		return false, err
	}
	return exist, nil
}

func (repo *ImgPSQLRepository) Create(img *domain.Img, categoryID int64) (*domain.Img, error) {
	ormImg := domainImgToORM(img)

	var category *domain.Category
	var err error
	// 如果有分类ID，在插入前设置
	if categoryID != 0 {
		ormImg.CategoryID.Valid = true
		ormImg.CategoryID.Int64 = categoryID

		category, err = repo.FindCategoryByID(img.TenantID, categoryID)
		if err != nil {
			return nil, err
		}
	}

	if category != nil {
		ormImg.Path = category.Prefix + "/" + ormImg.Path
	}

	if err := ormImg.InsertG(boil.Infer()); err != nil {
		return nil, err
	}

	return ormImgToDomain(ormImg), nil
}

func (repo *ImgPSQLRepository) Delete(tenantID domain.TenantID, id int64, hard bool) error {
	rows, err := orm.Imgs(
		qm.Where(fmt.Sprintf("%s = ?", orm.ImgColumns.TenantID), tenantID),
		qm.And(fmt.Sprintf("%s = ?", orm.ImgColumns.ID), id),
		qm.WithDeleted(),
	).DeleteAllG(hard)
	if err != nil {
		return err
	}
	if rows == 0 {
		return codes.ErrImgNotFound
	}
	return nil
}

func (repo *ImgPSQLRepository) Restore(tenantID domain.TenantID, id int64) (*domain.Img, error) {
	rows, err := orm.Imgs(
		qm.Where(fmt.Sprintf("%s = ?", orm.ImgColumns.TenantID), tenantID),
		qm.And(fmt.Sprintf("%s = ?", orm.ImgColumns.ID), id),
		qm.WithDeleted(),
	).UpdateAllG(orm.M{
		orm.ImgColumns.DeletedAt: nil,
		orm.ImgColumns.UpdatedAt: time.Now(),
	})
	if err != nil {
		return nil, err
	}

	if rows == 0 {
		return nil, codes.ErrImgNotFound
	}

	img, err := repo.FindByID(tenantID, id)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func (repo *ImgPSQLRepository) List(query *domain.ImgQuery) (*domain.ImgList, error) {
	var whereMods []qm.QueryMod

	whereMods = append(whereMods, qm.Where(fmt.Sprintf("%s = ?", orm.ImgColumns.TenantID), query.TenantID))

	if query.Keyword != "" {
		like := "%" + query.Keyword + "%"
		whereMods = append(whereMods, qm.Where("description ILIKE ?", like))
	}

	if query.CategoryID != 0 {
		whereMods = append(whereMods, qm.Where(fmt.Sprintf("%s = ?", orm.ImgColumns.CategoryID), query.CategoryID))
	}

	if query.Deleted {
		whereMods = append(whereMods, qm.WithDeleted())
		whereMods = append(whereMods, qm.Where("deleted_at is not null"))
	}

	// 1.计算count
	total, err := orm.Imgs(whereMods...).CountG()
	if err != nil {
		return nil, err
	}

	// 2.计算offset
	offset, err := dbkit.ComputeOffset(query.Page, query.PageSize)
	if err != nil {
		return nil, err
	}

	pageMods := append(whereMods, qm.Offset(offset), qm.Limit(query.PageSize), qm.OrderBy(orm.ImgColumns.UpdatedAt+" DESC"), qm.OrderBy(orm.ImgColumns.ID+" DESC"))

	imgs, err := orm.Imgs(pageMods...).AllG()
	if err != nil {
		return nil, err
	}

	return &domain.ImgList{
		Total: total,
		List:  ormImgsToDomain(imgs),
	}, nil
}

func (repo *ImgPSQLRepository) CreateCategory(category *domain.Category) error {
	ormCategory := domainCategoryToORM(category)
	return ormCategory.InsertG(boil.Infer())
}

func (repo *ImgPSQLRepository) UpdateCategory(category *domain.Category) error {
	ormCategory := domainCategoryToORM(category)
	rows, err := ormCategory.UpdateG(boil.Infer())
	if err != nil {
		return err
	}

	if rows == 0 {
		return codes.ErrImgNotFound
	}

	return nil
}

func (repo *ImgPSQLRepository) DeleteCategory(tenantID domain.TenantID, id int64) error {
	ormCategory := orm.ImgCategory{
		ID:       id,
		TenantID: int64(tenantID),
	}

	rows, err := ormCategory.DeleteG()
	if err != nil {
		return err
	}
	if rows == 0 {
		return codes.ErrImgCategoryNotFound
	}
	return nil
}

func (repo *ImgPSQLRepository) ListCategories(tenantID domain.TenantID) ([]*domain.Category, error) {
	ormCategories, err := orm.ImgCategories(
		qm.Where(fmt.Sprintf("%s = ?", orm.ImgCategoryColumns.TenantID), tenantID),
	).AllG()
	if err != nil {
		return nil, err
	}

	return ormCategoriesToDomain(ormCategories), nil
}

func (repo *ImgPSQLRepository) FindCategoryByID(tenantID domain.TenantID, id int64) (*domain.Category, error) {
	ormCategory, err := orm.ImgCategories(
		qm.Where(fmt.Sprintf("%s = ?", orm.ImgCategoryColumns.TenantID), tenantID),
		qm.And(fmt.Sprintf("%s = ?", orm.ImgCategoryColumns.ID), id),
	).OneG()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, codes.ErrImgCategoryNotFound
		}
		return nil, err
	}

	return ormCategoryToDomain(ormCategory), nil
}

func (repo *ImgPSQLRepository) FindCategoryByTitle(tenantID domain.TenantID, title string) (*domain.Category, error) {
	ormCategory, err := orm.ImgCategories(
		qm.Where(fmt.Sprintf("%s = ?", orm.ImgCategoryColumns.TenantID), tenantID),
		qm.And(fmt.Sprintf("%s = ?", orm.ImgCategoryColumns.Title), title),
	).OneG()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, codes.ErrImgCategoryNotFound
		}
		return nil, err
	}
	return ormCategoryToDomain(ormCategory), nil
}

func (repo *ImgPSQLRepository) CategoryExistByID(tenantID domain.TenantID, id int64) (bool, error) {
	exist, err := orm.Imgs(
		qm.Where(fmt.Sprintf("%s = ?", orm.ImgCategoryColumns.TenantID), tenantID),
		qm.And(fmt.Sprintf("%s = ?", orm.ImgColumns.CategoryID), id),
	).ExistsG()

	if err != nil {
		return false, err
	}

	return exist, nil
}

func (repo *ImgPSQLRepository) CategoryExistByTitle(tenantID domain.TenantID, title string) (bool, error) {
	exist, err := orm.ImgCategories(
		qm.Where(fmt.Sprintf("%s = ?", orm.ImgCategoryColumns.TenantID), tenantID),
		qm.And(fmt.Sprintf("%s = ?", orm.ImgCategoryColumns.Title), title),
	).ExistsG()
	if err != nil {
		return false, err
	}
	return exist, nil
}

func (repo *ImgPSQLRepository) CountCategory(tenantID domain.TenantID) (int64, error) {
	count, err := orm.ImgCategories(
		qm.Where(fmt.Sprintf("%s = ?", orm.ImgCategoryColumns.TenantID), tenantID),
	).CountG()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *ImgPSQLRepository) IsCategoryExistImg(tenantID domain.TenantID, id int64) (bool, error) {
	existing2, err := orm.Imgs(
		qm.Where(fmt.Sprintf("%s = ?", orm.ImgCategoryColumns.TenantID), tenantID),
		qm.And(fmt.Sprintf("%s = ?", orm.ImgColumns.CategoryID), id),
		qm.WithDeleted(),
	).ExistsG()
	if err != nil {
		return false, err
	}

	return existing2, nil
}

func (repo *ImgPSQLRepository) GetTenantR2Config(tenantID domain.TenantID) (*domain.R2Config, error) {
	config, err := orm.TenantR2Configs(
		qm.Where(fmt.Sprintf("%s = ?", orm.TenantR2ConfigColumns.TenantID), tenantID),
	).OneG()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, codes.ErrImgR2ConfigNotFound
		}
		return nil, err
	}

	return ormR2ConfigToDomian(config), nil
}

func (repo *ImgPSQLRepository) SetTenantR2Config(config *domain.R2Config) error {
	ormR2Config := doaminR2ConfigToORM(config)

	var updateWhiteList boil.Columns
	// 除开SecretAccessKey列 其余列必定不为空 故此仅需处理SecretAccessKey列
	if config.GetSecretAccessKey() == "" {
		updateWhiteList = boil.Blacklist(
			orm.TenantR2ConfigColumns.SecretAccessKey,
		)
	} else {
		updateWhiteList = boil.Infer()
	}

	err := ormR2Config.UpsertG(
		true,
		[]string{orm.TenantR2ConfigColumns.TenantID},
		updateWhiteList,
		boil.Infer(),
	)

	return err
}

func (repo *ImgPSQLRepository) ExistTenantR2Config(tenantID domain.TenantID) (bool, error) {
	exist, err := orm.TenantR2Configs(
		orm.TenantR2ConfigWhere.TenantID.EQ(int64(tenantID)),
	).ExistsG()
	if err != nil {
		return false, errors.WithStack(err)
	}

	return exist, nil
}
