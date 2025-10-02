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
}

type TenantQuery struct {
	Keyword  string
	Page     int
	PageSize int
}

type TenantList struct {
	Total int64
	List  []*Tenant
}
