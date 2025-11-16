package domain

import (
	"time"
)

type Tenant struct {
	ID          string
	PlanType    PlanType
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CreatorID   string
}

func (t *Tenant) GetCreatedAt() time.Time {
	return t.CreatedAt
}

func (t *Tenant) GetID() string {
	return t.ID
}

type TenantPagingQuery struct {
	PageSize   int
	CreatorID  string
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
