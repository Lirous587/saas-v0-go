package handler

type PlanResponse struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description,omitempty"`
	CreatedAt   int64   `json:"created_at"`
	UpdatedAt   int64   `json:"updated_at"`
}

type CreateRequest struct {
	Name        string  `json:"name" binding:"required,max=50"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
}

type UpdateRequest struct {
	ID          int64   `json:"-" uri:"id" binding:"required"`
	Name        string  `json:"name" binding:"required,max=50"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
}

type PlanListResponse struct {
	List []*PlanResponse `json:"list"`
}
