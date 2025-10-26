package domain

import (
	"time"
)

type Tenant struct {
	ID          int64
	PlanType    PlanType
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CreatorID   int64
}

func (t Tenant) GetCreatedAt() time.Time {
	return t.CreatedAt
}

func (t Tenant) GetID() int64 {
	return t.ID
}

type TenantPagingQuery struct {
	PageSize   int
	CreatorID  int64
	PrevCursor string
	NextCursor string
	Keyword    string
}

type TenantPagination struct {
	Items      []*Tenant
	PrevCursor string
	NextCursor string
	HasPrev    bool
	HasNext    bool
}
