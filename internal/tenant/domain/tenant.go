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

type Plan struct {
	ID   int64
	Name string
}

type TenantWithPlan struct {
	Tenant *Tenant
	Plan   *Plan
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
