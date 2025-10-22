package domain

import "database/sql"

type PlanRepository interface {
	Create(plan *Plan) (*Plan, error)
	Update(plan *Plan) error
	Delete(id int64) error
	List() ([]*Plan, error)
	CreatorHasPlan(creatorID, planID int64) (bool, error)

	AttchToTenantTx(tx *sql.Tx, planID, tenantID, creatorID int64) error
}

type PlanCache interface {
}
