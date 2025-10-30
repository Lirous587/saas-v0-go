package domain

type ImgRepository interface {
	FindByID(tenantID TenantID, id string, deleted ...bool) (*Img, error)
	ExistByPath(tenantID TenantID, path string) (bool, error)

	Create(img *Img, categoryID string) (*Img, error)
	Delete(tenantID TenantID, id string, hard bool) error
	Restore(tenantID TenantID, id string) (*Img, error)
	List(query *ImgQuery) (*ImgList, error)

	CreateCategory(category *Category) error
	UpdateCategory(category *Category) error
	DeleteCategory(tenantID TenantID, id string) error
	ListCategories(tenantID TenantID) ([]*Category, error)
	FindCategoryByID(tenantID TenantID, id string) (*Category, error)
	FindCategoryByTitle(tenantID TenantID, title string) (*Category, error)
	CategoryExistByTitle(tenantID TenantID, title string) (bool, error)
	CategoryExistByID(tenantID TenantID, id string) (bool, error)
	CountCategory(tenantID TenantID) (int64, error)
	IsCategoryExistImg(tenantID TenantID, id string) (bool, error)

	SetTenantR2Config(config *R2Config) error
	GetTenantR2Config(tenantID TenantID) (*R2Config, error)
	ExistTenantR2Config(tenantID TenantID) (bool, error)
}

type ImgMsgQueue interface {
	AddToDeleteQueue(tenantID TenantID, imgID string) error
	ListenDeleteQueue(onExpire func(tenantID TenantID, imgID string))
	RemoveFromDeleteQueue(tenantID TenantID, imgID string) error
}
