package domain

type CommentRepository interface {
	GetByID(tenantID TenantID, id int64) (*Comment, error)
	Create(comment *Comment) (*Comment, error)
	Delete(tenantID TenantID, id int64) error
	Approve(tenantID TenantID, id int64) error
	List(query *CommentQuery) (*CommentList, error)

	GetCommentUser(tenantID TenantID, commentID int64) (int64, error)

	GetUserIdsByRootORParent(tenantID TenantID, plateID int64, rootID int64, parentID int64) ([]int64, error)
	GetDomainAdminByTenant(tenantID TenantID) (*UserInfo, error)
	GetUserInfosByIds(ids []int64) ([]*UserInfo, error)
	GetUserInfoByID(id int64) (*UserInfo, error)

	SetTenantConfig(config *TenantConfig) error
	GetTenantConfig(tenantID TenantID) (*TenantConfig, error)
	ExistTenantConfigByID(tenantID TenantID) (bool, error)

	CreatePlate(plate *Plate) error
	DeletePlate(tenantID TenantID, id int64) error
	ListPlate(query *PlateQuery) (*PlateList, error)
	ExistPlateBykey(tenantID TenantID, belongKey string) (bool, error)
	GetPlateBelongByID(id int64) (*PlateBelong, error)
	GetPlateBelongByKey(tenantID TenantID, belongKey string) (*PlateBelong, error)
	GetPlateRelatedURlByID(tenantID TenantID, id int64) (string, error)
	SetPlateConfig(config *PlateConfig) error
	GetPlateConfig(tenantID TenantID, palteID int64) (*PlateConfig, error)
}

type CommentCache interface {
	SetTenantConfig(config *TenantConfig) error
	GetTenantConfig(tenantID TenantID) (*TenantConfig, error)
	DeleteTenantConfig(tenantID TenantID) error

	SetPlateConfig(config *PlateConfig) error
	GetPlateConfig(tenantID TenantID, plateID int64) (*PlateConfig, error)
	DeletePlateConfig(tenantID TenantID, plateID int64) error
}
