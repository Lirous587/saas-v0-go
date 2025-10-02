package domain

import (
	"time"
)

type Role struct {
	ID          int64
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type RoleQuery struct {
	Keyword  string
	Page     int
	PageSize int
}

type RoleList struct {
	Total int64
	List  []*Role
}
