package domain

import (
	"time"
)

type Tenant struct {
	ID          int64
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CreatorID   int64
}

type TenantPagingQuery struct {
	PageSize  int
	CreatorID int64
	AfterID   int64
	BeforeID  int64
	Keyword   string
}

type TenantPagination struct {
	Items   []*Tenant
	HasNext bool
	HasPrev bool
}
