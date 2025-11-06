package domain

import (
	"io"
)

type ImgService interface {
	Upload(src io.Reader, img *Img, categoryID CategoryID) error
	Delete(tenantID TenantID, imgID ImgID, hard ...bool) error
	List(query *ImgQuery) (*ImgList, error)
	ClearRecycleBin(tenantID TenantID, imgID ImgID) error
	ListenDeleteQueue()
	RestoreFromRecycleBin(tenantID TenantID, imgID ImgID) error

	//	分类
	CreateCategory(category *Category) error
	UpdateCategory(category *Category) error
	DeleteCategory(tenantID TenantID, categoryID CategoryID) error
	ListCategories(tenantID TenantID) (categories []*Category, err error)

	// 配置
	SetR2Config(config *R2Config) error
	GetR2Config(tenantID TenantID) (*R2Config, error)

	SetR2SecretKey(tenantID TenantID, secretKey R2SecretAccessKey) error
	IsSetR2SecretKey(tenantID TenantID) (bool, error)
}
