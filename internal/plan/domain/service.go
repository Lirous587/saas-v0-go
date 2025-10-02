package domain

type PlanService interface {
	Create(plan *Plan) (*Plan, error)
	Read(id int64) (*Plan, error)
	Update(plan *Plan) (*Plan, error)
	Delete(id int64) error
	List(query *PlanQuery) (*PlanList, error)
}
