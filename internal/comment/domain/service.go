package domain

type CommentService interface {
	Create(comment *Comment) (*Comment, error)
	Read(id int64) (*Comment, error)
	Update(comment *Comment) (*Comment, error)
	Delete(id int64) error
	List(query *CommentQuery) (*CommentList, error)
}
