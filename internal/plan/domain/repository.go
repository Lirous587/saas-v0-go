package domain


type PlanRepository interface {
	FindByID(id int64) (*Plan, error)

	Create(plan *Plan) (*Plan, error)
	Update(plan *Plan) (*Plan, error)
	Delete(id int64) error
	List(query *PlanQuery) (*PlanList, error)
}

type PlanCache interface {

}
