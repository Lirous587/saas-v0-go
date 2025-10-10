package domain

type CommentService interface {
	Create(comment *Comment) (*Comment, error)
	Read(id int64) (*Comment, error)
	Update(comment *Comment) (*Comment, error)
	Delete(id int64) error
	List(query *CommentQuery) (*CommentList, error)

	CreatePlate(plate *Plate) error
	DeletePlate(tenantID TenantID, id int64) error
	ListPlate(query *PlateQuery) (*PlateList, error)

	SetTenantConfig(config *TenantConfig) error
	GetTenantConfig(tenantID TenantID) (*TenantConfig, error)
	SetPlateConfig(config *PlateConfig) error
	GetPlateConfig(tenantID TenantID, plate string) (*PlateConfig, error)
}
