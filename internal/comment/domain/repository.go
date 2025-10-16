package domain

type CommentRepository interface {
	GetByID(tenantID TenantID, id int64) (*Comment, error)
	Create(comment *Comment) (*Comment, error)
	Delete(tenantID TenantID, id int64) error
	Approve(tenantID TenantID, id int64) error
	ListRoots(query *CommentRootsQuery) ([]*CommentRoot, error)
	ListReplies(query *CommentRepliesQuery) ([]*CommentReply, error)
	UpdateLikeCount(tenantID TenantID, commentID int64, isLike bool) error

	GetCommentUser(tenantID TenantID, commentID int64) (int64, error)

	GetUserIdsByRootORParent(tenantID TenantID, plateID int64, rootID int64, parentID int64) ([]int64, error)
	GetDomainAdminByTenant(tenantID TenantID) (*UserInfo, error)
	GetUserInfosByIds(ids []int64) ([]*UserInfo, error)
	GetUserInfoByID(id int64) (*UserInfo, error)

	SetTenantConfig(config *TenantConfig) error
	GetTenantConfig(tenantID TenantID) (*TenantConfig, error)
	ExistTenantConfigByID(tenantID TenantID) (bool, error)

	CreatePlate(plate *Plate) error
	UpdatePlate(plate *Plate) error
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

	GetPlateID(tenantID TenantID, belongKey string) (int64, error)
	SetPlateID(tenantID TenantID, belongKey string, id int64) error
	DeletePlateID(tenantID TenantID, belongKey string) error
	SetPlateConfig(config *PlateConfig) error
	GetPlateConfig(tenantID TenantID, plateID int64) (*PlateConfig, error)
	DeletePlateConfig(tenantID TenantID, plateID int64) error

	GetLikeStatus(tenantID TenantID, userID int64, commentID int64) (LikeStatus, error)
	AddLike(tenantID TenantID, userID int64, commentID int64) error
	RemoveLike(tenantID TenantID, userID int64, commentID int64) error
	GetLikeMap(tenantID TenantID, userID int64, commentIds []int64) (map[int64]struct{}, error)
}
