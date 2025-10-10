package domain

type CommentRepository interface {
	FindByID(id int64) (*Comment, error)
	Create(comment *Comment) (*Comment, error)
	Update(comment *Comment) (*Comment, error)
	Delete(id int64) error
	List(query *CommentQuery) (*CommentList, error)

	CreatePlate(plate *Plate) error
	DeletePlate(tenantID TenantID, id int64) error
	ExistPlate(tenantID TenantID, palte string) (bool, error)
	ListPlate(query *PlateQuery) (*PlateList, error)

	SetTenantConfig(config *TenantConfig) error
	GetTenantConfig(tenantID TenantID) (*TenantConfig, error)
	SetPlateConfig(config *PlateConfig) error
	GetPlateConfig(tenantID TenantID, palte string) (*PlateConfig, error)
}

type CommentCache interface {
	SetTenantCommentClientToken(tenantID TenantID, clientToken string) error
	GetTenantCommentClientToken(tenantID TenantID) (string, error)
}
