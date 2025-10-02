package domain

import (
	"time"
)

type Plan struct {
	ID          int64
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type PlanQuery struct {
	Keyword  string
	Page     int
	PageSize int
}

type PlanList struct {
	Total int64
	List  []*Plan
}
