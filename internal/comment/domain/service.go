package domain

type CommentService interface {
	Create(comment *Comment, belongKey string) error
	Delete(tenantID TenantID, userID string, id string) error
	Audit(tenantID TenantID, id string, status CommentStatus) error
	ListRoots(belongKey string, userID string, query *CommentRootsQuery) ([]*CommentRoot, error)
	ListReplies(belongKey string, userID string, query *CommentRepliesQuery) ([]*CommentReply, error)

	ToggleLike(tenantID TenantID, userID string, id string) error
	CreatePlate(plate *Plate) error
	UpdatePlate(plate *Plate) error
	DeletePlate(tenantID TenantID, id string) error
	ListPlate(query *PlateQuery) (*PlateList, error)

	SetTenantConfig(config *TenantConfig) error
	GetTenantConfig(tenantID TenantID) (*TenantConfig, error)
	SetPlateConfig(config *PlateConfig) error
	GetPlateConfig(tenantID TenantID, plateID string) (*PlateConfig, error)
}
