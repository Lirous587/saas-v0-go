package domain

type ImgRepository interface {
	FindByID(tenantID TenantID, id int64, deleted ...bool) (*Img, error)
	ExistByPath(tenantID TenantID, path string) (bool, error)

	Create(img *Img, categoryID int64) (*Img, error)
	Delete(tenantID TenantID, id int64, hard bool) error
	Restore(tenantID TenantID, id int64) (*Img, error)
	List(query *ImgQuery) (*ImgList, error)

	CreateCategory(category *Category) error
	UpdateCategory(category *Category) error
	DeleteCategory(tenantID TenantID, id int64) error
	ListCategories(tenantID TenantID) ([]*Category, error)
	FindCategoryByID(tenantID TenantID, id int64) (*Category, error)
	FindCategoryByTitle(tenantID TenantID, title string) (*Category, error)
	CategoryExistByTitle(tenantID TenantID, title string) (bool, error)
	CategoryExistByID(tenantID TenantID, id int64) (bool, error)
	CountCategory(tenantID TenantID) (int64, error)
	IsCategoryExistImg(tenantID TenantID, id int64) (bool, error)

	SetTenantR2Config(config *R2Config) error
	GetTenantR2Config(tenantID TenantID) (*R2Config, error)
	ExistTenantR2Config(tenantID TenantID) (bool, error)
}

type ImgMsgQueue interface {
	AddToDeleteQueue(tenantID TenantID, imgID int64) error
	ListenDeleteQueue(onExpire func(tenantID TenantID, imgID int64))
	RemoveFromDeleteQueue(tenantID TenantID, imgID int64) error
}
