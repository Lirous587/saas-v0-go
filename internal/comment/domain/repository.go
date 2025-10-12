package domain

type CommentRepository interface {
	// FindByID(id int64) (*Comment, error)
	Create(comment *Comment) (*Comment, error)
	Delete(tenantID TenantID, id int64) error
	List(query *CommentQuery) (*CommentList, error)
	IsCommentInPlate(tenantID TenantID, plateID int64, commentID int64) (bool, error)

	GetUserIdsByRootORParent(tenantID TenantID, plateID int64, rootID int64, parentID int64) ([]int64, error)
	GetDomainAdminByTenant(tenantID TenantID) (int64, error)
	GetUserInfosByIds(ids []int64) ([]*UserInfo, error)
	GetUserInfoByID(id int64) (*UserInfo, error)

	CreatePlate(plate *Plate) error
	DeletePlate(tenantID TenantID, id int64) error
	ListPlate(query *PlateQuery) (*PlateList, error)
	ExistPlateBykey(tenantID TenantID, belongKey string) (bool, error)
	GetPlateBelongByID(id int64) (*PlateBelong, error)
	GetPlateBelongByKey(tenantID TenantID, belongKey string) (*PlateBelong, error)
	GetPlateRelatedURlByID(tenantID TenantID, id int64) (string, error)

	SetTenantConfig(config *TenantConfig) error
	GetTenantConfig(tenantID TenantID) (*TenantConfig, error)
	SetPlateConfig(config *PlateConfig) error
	GetPlateConfig(tenantID TenantID, palteID int64) (*PlateConfig, error)
}

type CommentCache interface {
	SetTenantCommentClientToken(tenantID TenantID, clientToken string) error
	GetTenantCommentClientToken(tenantID TenantID) (string, error)
}
