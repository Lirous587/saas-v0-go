package domain


type CommentRepository interface {
	FindByID(id int64) (*Comment, error)

	Create(comment *Comment) (*Comment, error)
	Update(comment *Comment) (*Comment, error)
	Delete(id int64) error
	List(query *CommentQuery) (*CommentList, error)
}

type CommentCache interface {

}
