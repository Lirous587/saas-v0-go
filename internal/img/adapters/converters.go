package adapters

import (
	"github.com/aarondl/null/v8"
	"saas/internal/common/orm"
	"saas/internal/img/domain"
)

func domainImgToORM(img *domain.Img) *orm.Img {
	if img == nil {
		return nil
	}

	ormImg := &orm.Img{
		ID:        img.ID,
		TenantID:  int64(img.TenantID),
		Path:      img.Path,
		UpdatedAt: img.UpdatedAt,
	}

	// 处理null项
	if img.Description != "" {
		ormImg.Description = null.StringFrom(img.Description)
	}

	return ormImg
}

func ormImgToDomain(ormImg *orm.Img, isDeleted ...bool) *domain.Img {
	if ormImg == nil {
		return nil
	}

	img := &domain.Img{
		ID:        ormImg.ID,
		TenantID:  domain.TenantID(ormImg.TenantID),
		Path:      ormImg.Path,
		CreatedAt: ormImg.CreatedAt,
		UpdatedAt: ormImg.UpdatedAt,
	}

	// 处理null项
	if ormImg.Description.Valid {
		img.Description = ormImg.Description.String
	}

	if ormImg.DeletedAt.Valid {
		img.DeletedAt = ormImg.DeletedAt.Time
	}

	return img
}

func ormImgsToDomain(ormImgs []*orm.Img, isDeleted ...bool) []*domain.Img {
	if len(ormImgs) == 0 {
		return nil
	}

	deleted := false
	if len(isDeleted) > 0 {
		deleted = isDeleted[0]
	}

	imgs := make([]*domain.Img, 0, len(ormImgs))
	for _, ormImg := range ormImgs {
		if ormImg == nil {
			continue
		}
		img := ormImgToDomain(ormImg, deleted)
		imgs = append(imgs, img)
	}

	return imgs
}

func domainCategoryToORM(category *domain.Category) *orm.ImgCategory {
	if category == nil {
		return nil
	}

	ormImg := &orm.ImgCategory{
		ID:       category.ID,
		TenantID: int64(category.TenantID),
		Title:    category.Title,
		Prefix:   category.Prefix,
	}

	// 处理null项

	return ormImg
}

func ormCategoryToDomain(ormCategory *orm.ImgCategory) *domain.Category {
	if ormCategory == nil {
		return nil
	}

	img := &domain.Category{
		ID:        ormCategory.ID,
		TenantID:  domain.TenantID(ormCategory.TenantID),
		Title:     ormCategory.Title,
		Prefix:    ormCategory.Prefix,
		CreatedAt: ormCategory.CreatedAt,
	}

	// 处理null项
	return img
}

func ormCategoriesToDomain(ormCategories []*orm.ImgCategory) []*domain.Category {
	if len(ormCategories) == 0 {
		return nil
	}

	categories := make([]*domain.Category, 0, len(ormCategories))
	for _, ormCategory := range ormCategories {
		if ormCategory == nil {
			continue
		}
		category := ormCategoryToDomain(ormCategory)
		categories = append(categories, category)
	}

	return categories
}

func doaminR2ConfigToORM(r2Config *domain.R2Config) *orm.TenantR2Config {
	if r2Config == nil {
		return nil
	}

	ormR2Config := &orm.TenantR2Config{
		TenantID:        int64(r2Config.TenantID),
		AccountID:       r2Config.AccountID,
		AccessKeyID:     r2Config.AccessKeyID,
		SecretAccessKey: r2Config.GetSecretAccessKey(),
		PublicBucket:    r2Config.PublicBucket,
		PublicURLPrefix: r2Config.PublicURLPrefix,
		DeleteBucket:    r2Config.DeleteBucket,
	}

	return ormR2Config
}

func ormR2ConfigToDomian(ormR2Config *orm.TenantR2Config) *domain.R2Config {
	if ormR2Config == nil {
		return nil
	}

	r2Config := &domain.R2Config{
		TenantID:        domain.TenantID(ormR2Config.TenantID),
		AccountID:       ormR2Config.AccountID,
		AccessKeyID:     ormR2Config.AccessKeyID,
		PublicBucket:    ormR2Config.PublicBucket,
		PublicURLPrefix: ormR2Config.PublicURLPrefix,
		DeleteBucket:    ormR2Config.DeleteBucket,
	}

	r2Config.SetSecretAccessKey(ormR2Config.SecretAccessKey)

	return r2Config
}
