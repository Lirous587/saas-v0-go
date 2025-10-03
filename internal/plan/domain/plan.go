package domain

import (
	"time"
)

type Plan struct {
	ID          int64
	Name        string
	Price       float64
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type PlanList struct {
	List []*Plan
}
