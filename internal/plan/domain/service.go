package domain

import "database/sql"

type PlanService interface {
	Create(plan *Plan) error
	Update(plan *Plan) error
	Delete(id int64) error
	List() ([]*Plan, error)
	CreatorHasPlan(creatorID, planID int64) (bool, error)

	AttchToTenantTx(tx *sql.Tx, planID, tenantID, creatorID int64) error
}
