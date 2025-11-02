package domain

type ImgRepository interface {
	FindByID(tenantID TenantID, imgID ImgID, deleted ...bool) (*Img, error)
	ExistByPath(tenantID TenantID, path string) (bool, error)

	Create(img *Img, categoryID CategoryID) (*Img, error)
	Delete(tenantID TenantID, imgID ImgID, hard bool) error
	Restore(tenantID TenantID, imgID ImgID) (*Img, error)
	List(query *ImgQuery) (*ImgList, error)

	CreateCategory(category *Category) error
	UpdateCategory(category *Category) error
	DeleteCategory(tenantID TenantID, categoryID CategoryID) error
	ListCategories(tenantID TenantID) ([]*Category, error)
	FindCategoryByID(tenantID TenantID, categoryID CategoryID) (*Category, error)
	FindCategoryByTitle(tenantID TenantID, title string) (*Category, error)
	CategoryExistByTitle(tenantID TenantID, title string) (bool, error)
	CategoryExistByID(tenantID TenantID, categoryID CategoryID) (bool, error)
	CountCategory(tenantID TenantID) (int64, error)
	IsCategoryExistImg(tenantID TenantID, categoryID CategoryID) (bool, error)

	SetTenantR2Config(config *R2Config) error
	GetTenantR2Config(tenantID TenantID) (*R2Config, error)
	ExistTenantR2Config(tenantID TenantID) (bool, error)
}

type ImgMsgQueue interface {
	AddToDeleteQueue(tenantID TenantID, imgID ImgID) error
	ListenDeleteQueue(onExpire func(tenantID TenantID, imgID ImgID))
	RemoveFromDeleteQueue(tenantID TenantID, imgID ImgID) error
}
