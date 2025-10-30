package domain

import (
	"io"
)

type ImgService interface {
	Upload(src io.Reader, img *Img, categoryID string) error
	Delete(tenantID TenantID, id string, hard ...bool) error
	List(query *ImgQuery) (*ImgList, error)
	ClearRecycleBin(tenantID TenantID, id string) error
	ListenDeleteQueue()
	RestoreFromRecycleBin(tenantID TenantID, id string) error

	//	分类
	CreateCategory(category *Category) error
	UpdateCategory(category *Category) error
	DeleteCategory(tenantID TenantID, id string) error
	ListCategories(tenantID TenantID) (categories []*Category, err error)

	// 配置
	SetR2Config(config *R2Config) error
	GetR2Config(tenantID TenantID) (*R2Config, error)
}
