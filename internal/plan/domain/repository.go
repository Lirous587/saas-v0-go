package domain

import "database/sql"

type PlanRepository interface {
	FindByID(id int64) (*Plan, error)

	Create(plan *Plan) (*Plan, error)
	Update(plan *Plan) (*Plan, error)
	Delete(id int64) error
	List() (*PlanList, error)

	AttchToTenantTx(tx *sql.Tx, planID, tenantID int64) error
}

type PlanCache interface {
}
