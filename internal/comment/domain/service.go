package domain

type CommentService interface {
	Create(comment *Comment, belongKey string) error
	Delete(tenantID TenantID, userID UserID, commentID CommentID) error
	Audit(tenantID TenantID, commentID CommentID, status CommentStatus) error
	ListRoots(belongKey string, userID UserID, query *CommentRootsQuery) ([]*CommentRoot, error)
	ListReplies(belongKey string, userID UserID, query *CommentRepliesQuery) ([]*CommentReply, error)

	ToggleLike(tenantID TenantID, userID UserID, commentID CommentID) error

	CreatePlate(plate *Plate) error
	UpdatePlate(plate *Plate) error
	DeletePlate(tenantID TenantID, plateID PlateID) error
	ListPlate(query *PlateQuery) (*PlateList, error)
	CheckPlateBelongKey(tenantID TenantID,belongKey string) (bool,error)

	SetTenantConfig(config *TenantConfig) error
	GetTenantConfig(tenantID TenantID) (*TenantConfig, error)
	SetPlateConfig(config *PlateConfig) error
	GetPlateConfig(tenantID TenantID, plateID PlateID) (*PlateConfig, error)
}
