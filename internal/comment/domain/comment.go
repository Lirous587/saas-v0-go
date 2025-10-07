package domain

import (
	"time"
)

type Comment struct {
	ID          int64
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CommentQuery struct {
	Keyword  string
	Page     int
	PageSize int
}

type CommentList struct {
	Total int64
	List  []*Comment
}
