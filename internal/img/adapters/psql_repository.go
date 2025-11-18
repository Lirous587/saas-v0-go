package adapters

import (
	"database/sql"
	"fmt"
	"saas/internal/common/orm"
	"saas/internal/common/reskit/codes"
	"saas/internal/common/utils/dbkit"
	"saas/internal/img/domain"
	"time"

	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/pkg/errors"
)

type ImgPSQLRepository struct {
}

func NewImgPSQLRepository() domain.ImgRepository {
	return &ImgPSQLRepository{}
}

func (repo *ImgPSQLRepository) FindByID(tenantID domain.TenantID, imgID domain.ImgID, deleted ...bool) (*domain.Img, error) {
	selectDeleted := len(deleted) > 0 && deleted[0]

	var whereMods []qm.QueryMod
	whereMods = append(whereMods,
		orm.ImgWhere.TenantID.EQ(tenantID.String()),
		orm.ImgWhere.ID.EQ(imgID.String()),
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
		orm.ImgWhere.TenantID.EQ(tenantID.String()),
		orm.ImgWhere.Path.EQ(path),
		qm.WithDeleted(),
	).ExistsG()
	if err != nil {
		return false, err
	}
	return exist, nil
}

func (repo *ImgPSQLRepository) Create(img *domain.Img, categoryID domain.CategoryID) (*domain.Img, error) {
	ormImg := domainImgToORM(img)

	var category *domain.Category
	var err error
	// 如果有分类ID，在插入前设置
	if categoryID != "" {
		ormImg.CategoryID.Valid = true
		ormImg.CategoryID.String = categoryID.String()

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

func (repo *ImgPSQLRepository) Delete(tenantID domain.TenantID, imgID domain.ImgID, hard bool) error {
	rows, err := orm.Imgs(
		orm.ImgWhere.TenantID.EQ(tenantID.String()),
		orm.ImgWhere.ID.EQ(imgID.String()),
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

func (repo *ImgPSQLRepository) Restore(tenantID domain.TenantID, imgID domain.ImgID) (*domain.Img, error) {
	rows, err := orm.Imgs(
		orm.ImgWhere.TenantID.EQ(tenantID.String()),
		orm.ImgWhere.ID.EQ(imgID.String()),
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

	img, err := repo.FindByID(tenantID, imgID)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func (repo *ImgPSQLRepository) ListByKeyset(query *domain.ListByKeysetQuery) (*domain.ListByKeysetResult, error) {
	var baseMods []qm.QueryMod

	baseMods = append(
		baseMods,
		orm.ImgWhere.TenantID.EQ(query.TenantID.String()),
	)

	if query.Keyword != "" {
		like := "%" + query.Keyword + "%"
		baseMods = append(baseMods, qm.Where("description ILIKE ?", like))
	}

	if query.CategoryID != "" {
		baseMods = append(
			baseMods,
			qm.Where(fmt.Sprintf("%s = ?", orm.ImgColumns.CategoryID), query.CategoryID))
	}

	if query.Deleted {
		baseMods = append(baseMods, qm.WithDeleted())
		baseMods = append(baseMods, qm.Where("deleted_at is not null"))
	}

	// Keyset: 主排序为 created_at, tie-breaker 为 id
	ks := dbkit.NewKeyset(
		orm.ImgColumns.ID,
		orm.ImgColumns.CreatedAt,
		query.PrevCursor,
		query.NextCursor,
		query.PageSize,
		dbkit.WithPrimaryOrder[*domain.Img](dbkit.SortDesc), // 保持与原来的 CreatedAt DESC 一致
	)

	// 使用 keyset 生成包含 ORDER BY / LIMIT 的 query mods
	mods := ks.ApplyKeysetMods(baseMods)

	ormImgs, err := orm.Imgs(mods...).AllG()
	if err != nil {
		return nil, err
	}

	domains := ormImgsToDomain(ormImgs)

	// 精确判断 hasPrev/hasNext：exists 必须和 baseMods 保持一致
	exists := func(primary time.Time, id string, checkPrev bool) (bool, error) {
		var cond qm.QueryMod
		if checkPrev {
			cond = ks.BeforeWhere(primary, id)
		} else {
			cond = ks.AfterWhere(primary, id)
		}
		checkMods := append([]qm.QueryMod{}, baseMods...)
		checkMods = append(checkMods, cond, qm.Limit(1))
		return orm.Imgs(checkMods...).ExistsG()
	}

	// 精确构建分页结果（包含 HasPrev/HasNext, 游标）
	pager, err := ks.BuildPaginationResultWithExistence(domains, exists)
	if err != nil {
		return nil, err
	}

	return &domain.ListByKeysetResult{
		Items:      pager.Items,
		PrevCursor: pager.PrevCursor,
		NextCursor: pager.NextCursor,
		HasPrev:    pager.HasPrev,
		HasNext:    pager.HasNext,
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

func (repo *ImgPSQLRepository) DeleteCategory(tenantID domain.TenantID, categoryID domain.CategoryID) error {
	ormCategory := orm.ImgCategory{
		ID:       categoryID.String(),
		TenantID: tenantID.String(),
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

func (repo *ImgPSQLRepository) AllCategories(tenantID domain.TenantID) ([]*domain.Category, error) {
	ormCategories, err := orm.ImgCategories(
		orm.ImgCategoryWhere.TenantID.EQ(tenantID.String()),
	).AllG()
	if err != nil {
		return nil, err
	}

	return ormCategoriesToDomain(ormCategories), nil
}

func (repo *ImgPSQLRepository) FindCategoryByID(tenantID domain.TenantID, categoryID domain.CategoryID) (*domain.Category, error) {
	ormCategory, err := orm.ImgCategories(
		orm.ImgCategoryWhere.TenantID.EQ(tenantID.String()),
		qm.And(fmt.Sprintf("%s = ?", orm.ImgCategoryColumns.ID), categoryID),
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
		orm.ImgCategoryWhere.TenantID.EQ(tenantID.String()),
		orm.ImgCategoryWhere.Title.EQ(title),
	).OneG()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, codes.ErrImgCategoryNotFound
		}
		return nil, err
	}
	return ormCategoryToDomain(ormCategory), nil
}

func (repo *ImgPSQLRepository) CategoryExistByID(tenantID domain.TenantID, categoryID domain.CategoryID) (bool, error) {
	exist, err := orm.Imgs(
		orm.ImgCategoryWhere.TenantID.EQ(tenantID.String()),
		orm.ImgCategoryWhere.ID.EQ(categoryID.String()),
	).ExistsG()

	if err != nil {
		return false, err
	}

	return exist, nil
}

func (repo *ImgPSQLRepository) CategoryExistByTitle(tenantID domain.TenantID, title string) (bool, error) {
	exist, err := orm.ImgCategories(
		orm.ImgCategoryWhere.TenantID.EQ(tenantID.String()),
		orm.ImgCategoryWhere.Title.EQ(title),
	).ExistsG()
	if err != nil {
		return false, err
	}
	return exist, nil
}

func (repo *ImgPSQLRepository) CountCategory(tenantID domain.TenantID) (int64, error) {
	count, err := orm.ImgCategories(
		orm.ImgCategoryWhere.TenantID.EQ(tenantID.String()),
	).CountG()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *ImgPSQLRepository) IsCategoryExistImg(tenantID domain.TenantID, categoryID domain.CategoryID) (bool, error) {
	exist, err := orm.Imgs(
		orm.ImgWhere.TenantID.EQ(tenantID.String()),
		orm.ImgWhere.CategoryID.EQ(null.StringFrom(categoryID.String())),
		qm.WithDeleted(),
	).ExistsG()
	if err != nil {
		return false, err
	}

	return exist, nil
}

func (repo *ImgPSQLRepository) ExistTenantR2Config(tenantID domain.TenantID) (bool, error) {
	exist, err := orm.TenantR2Configs(
		orm.TenantR2ConfigWhere.TenantID.EQ(tenantID.String()),
	).ExistsG()
	if err != nil {
		return false, errors.WithStack(err)
	}

	return exist, nil
}

func (repo *ImgPSQLRepository) GetTenantR2Config(tenantID domain.TenantID) (*domain.R2Config, error) {
	config, err := orm.TenantR2Configs(
		orm.TenantR2ConfigWhere.TenantID.EQ(tenantID.String()),
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

	err := ormR2Config.UpsertG(
		true,
		[]string{orm.TenantR2ConfigColumns.TenantID},
		boil.Blacklist(
			orm.TenantR2ConfigColumns.SecretAccessKey,
		),
		boil.Blacklist(
			orm.TenantR2ConfigColumns.SecretAccessKey,
		),
	)

	return err
}

func (repo *ImgPSQLRepository) SetR2SecretKey(tenantID domain.TenantID, secretKey domain.R2SecretAccessKey) error {
	_, err := orm.TenantR2Configs(
		orm.TenantR2ConfigWhere.TenantID.EQ(string(tenantID)),
	).UpdateAllG(
		orm.M{
			orm.TenantR2ConfigColumns.SecretAccessKey: string(secretKey),
		},
	)

	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (repo *ImgPSQLRepository) IsSetR2SecretKey(tenantID domain.TenantID) (bool, error) {
	exist, err := orm.TenantR2Configs(
		orm.TenantR2ConfigWhere.TenantID.EQ(tenantID.String()),
		orm.TenantR2ConfigWhere.SecretAccessKey.IsNotNull(),
	).ExistsG()
	if err != nil {
		return false, errors.WithStack(err)
	}

	return exist, nil
}
