package domain

type CommentRepository interface {
	FindByID(id int64) (*Comment, error)

	Create(comment *Comment) (*Comment, error)
	Update(comment *Comment) (*Comment, error)
	Delete(id int64) error
	List(query *CommentQuery) (*CommentList, error)

	SetCommentTenantConfig(config *CommentTenantConfig) error
	GetCommentTenantConfig(tenantID TenantID) (*CommentTenantConfig, error)
	SetCommentConfig(config *CommentConfig) error
	GetCommentConfig(tenantID TenantID, benlongKey BelongKey) (*CommentConfig, error)
}

type CommentCache interface {
	SetTenantCommentClientToken(tenantID TenantID, clientToken string) error
	GetTenantCommentClientToken(tenantID TenantID, benlongKey BelongKey) (string, error)

	SetCommentClientToken(tenantID TenantID, belongKey BelongKey, clientToken string) error
	GetCommentClientToken(tenantID TenantID, benlongKey BelongKey) (string, error)
}
