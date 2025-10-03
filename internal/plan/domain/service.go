package domain

import "database/sql"

type PlanService interface {
	Create(plan *Plan) (*Plan, error)
	Read(id int64) (*Plan, error)
	Update(plan *Plan) (*Plan, error)
	Delete(id int64) error
	List() (*PlanList, error)

	AttchToTenantTx(tx *sql.Tx, planID, tenantID int64) error
}
