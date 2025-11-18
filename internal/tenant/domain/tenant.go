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

func (t *Tenant) GetCursorPrimary() time.Time {
	return t.UpdatedAt
}

func (t *Tenant) GetID() string {
	return t.ID
}

type TenantKeysetQuery struct {
	PageSize   int
	CreatorID  string
	PrevCursor string
	NextCursor string
	Keyword    string
}

type TenantKeysetResult struct {
	Items      []*Tenant
	PrevCursor string
	NextCursor string
	HasPrev    bool
	HasNext    bool
}
