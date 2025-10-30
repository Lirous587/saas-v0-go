package domain

type CommentRepository interface {
	GetByID(tenantID TenantID, id string) (*Comment, error)
	Create(comment *Comment) (*Comment, error)
	Delete(tenantID TenantID, id string) error
	Approve(tenantID TenantID, id string) error
	ListRoots(query *CommentRootsQuery) ([]*CommentRoot, error)
	ListReplies(query *CommentRepliesQuery) ([]*CommentReply, error)
	UpdateLikeCount(tenantID TenantID, commentID string, isLike bool) error

	GetCommentUser(tenantID TenantID, commentID string) (string, error)

	GetUserIdsByRootORParent(tenantID TenantID, plateID string, rootID string, parentID string) ([]string, error)
	GetTenantCreator(tenantID TenantID) (*UserInfo, error)
	GetUserInfosByIds(ids []string) ([]*UserInfo, error)
	GetUserInfoByID(id string) (*UserInfo, error)

	SetTenantConfig(config *TenantConfig) error
	GetTenantConfig(tenantID TenantID) (*TenantConfig, error)
	ExistTenantConfigByID(tenantID TenantID) (bool, error)

	CreatePlate(plate *Plate) error
	UpdatePlate(plate *Plate) error
	DeletePlate(tenantID TenantID, id string) error
	ListPlate(query *PlateQuery) (*PlateList, error)
	ExistPlateBykey(tenantID TenantID, belongKey string) (bool, error)
	GetPlateBelongByID(id string) (*PlateBelong, error)
	GetPlateBelongByKey(tenantID TenantID, belongKey string) (*PlateBelong, error)
	GetPlateRelatedURlByID(tenantID TenantID, id string) (string, error)
	SetPlateConfig(config *PlateConfig) error
	GetPlateConfig(tenantID TenantID, palteID string) (*PlateConfig, error)
}

type CommentCache interface {
	SetTenantConfig(config *TenantConfig) error
	GetTenantConfig(tenantID TenantID) (*TenantConfig, error)
	DeleteTenantConfig(tenantID TenantID) error

	GetPlateID(tenantID TenantID, belongKey string) (string, error)
	SetPlateID(tenantID TenantID, belongKey string, id string) error
	DeletePlateID(tenantID TenantID, belongKey string) error
	SetPlateConfig(config *PlateConfig) error
	GetPlateConfig(tenantID TenantID, plateID string) (*PlateConfig, error)
	DeletePlateConfig(tenantID TenantID, plateID string) error

	GetLikeStatus(tenantID TenantID, userID string, commentID string) (LikeStatus, error)
	AddLike(tenantID TenantID, userID string, commentID string) error
	RemoveLike(tenantID TenantID, userID string, commentID string) error
	GetLikeMap(tenantID TenantID, userID string, commentIds []string) (map[string]struct{}, error)
}
