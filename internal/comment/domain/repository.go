package domain

type CommentRepository interface {
	GetByID(tenantID TenantID, commentID CommentID) (*Comment, error)
	Create(comment *Comment) (*Comment, error)
	Delete(tenantID TenantID, commentID CommentID) error
	Approve(tenantID TenantID, commentID CommentID) error
	ListRoots(query *CommentRootsQuery) ([]*CommentRoot, error)
	ListReplies(query *CommentRepliesQuery) ([]*CommentReply, error)
	ListNoAudits(query *CommentNoAuditQuery) ([]*CommentNoAudit, error)
	UpdateLikeCount(tenantID TenantID, commentID CommentID, isLike bool) error

	GetCommentUser(tenantID TenantID, commentID CommentID) (UserID, error)

	GetUserIDsByRootORParent(tenantID TenantID, plateID PlateID, rootID CommentID, parentID CommentID) ([]UserID, error)
	GetTenantCreator(tenantID TenantID) (*UserInfo, error)
	GetUserInfosByIDs(userIDs []UserID) ([]*UserInfo, error)
	GetUserInfoByID(userID UserID) (*UserInfo, error)

	SetTenantConfig(config *TenantConfig) error
	GetTenantConfig(tenantID TenantID) (*TenantConfig, error)
	ExistTenantConfigByID(tenantID TenantID) (bool, error)

	CreatePlate(plate *Plate) error
	UpdatePlate(plate *Plate) error
	DeletePlate(tenantID TenantID, plateID PlateID) error
	ListPlate(query *PlateQuery) (*PlateList, error)
	ExistPlateBykey(tenantID TenantID, belongKey string) (bool, error)
	GetPlateBelongByID(plateID PlateID) (*PlateBelong, error)
	GetPlateBelongByKey(tenantID TenantID, belongKey string) (*PlateBelong, error)
	GetPlateRelatedURlByID(tenantID TenantID, plateID PlateID) (string, error)
	SetPlateConfig(config *PlateConfig) error
	GetPlateConfig(tenantID TenantID, palteID PlateID) (*PlateConfig, error)

	GetLikeStatus(tenantID TenantID, commentID CommentID, userID UserID) (LikeStatus, error)
	AddLike(tenantID TenantID, commentID CommentID, userID UserID) error
	RemoveLike(tenantID TenantID, commentID CommentID, userID UserID) error
	// GetLikeRecords 获取到点赞过的评论列表
	GetLikeRecords(tenantID TenantID, commentIDs []CommentID, userID UserID) ([]CommentID, error)
}

type CommentCache interface {
	SetTenantConfig(config *TenantConfig) error
	GetTenantConfig(tenantID TenantID) (*TenantConfig, error)
	DeleteTenantConfig(tenantID TenantID) error

	GetPlateID(tenantID TenantID, belongKey string) (PlateID, error)
	SetPlateID(tenantID TenantID, belongKey string, plateID PlateID) error
	DeletePlateID(tenantID TenantID, belongKey string) error
	SetPlateConfig(config *PlateConfig) error
	GetPlateConfig(tenantID TenantID, plateID PlateID) (*PlateConfig, error)
	DeletePlateConfig(tenantID TenantID, plateID PlateID) error

	// GetLikeMap(tenantID TenantID, userID UserID, commentIDs []CommentID) (map[CommentID]struct{}, error)
}
