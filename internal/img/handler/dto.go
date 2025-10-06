package handler

import "saas/internal/img/domain"

type ImgResponse struct {
	ID          int64  `json:"id"`
	Url         string `json:"url"`
	Description string `json:"description,omitempty"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

type UploadRequest struct {
	TenantID    domain.TenantID `json:"-" uri:"tenant_id" binding:"required"`
	Path        string          `form:"path" binding:"omitempty,slug"`
	Description string          `form:"description" binding:"max=60"`
	CategoryID  int64           `form:"category_id"`
}

type DeleteRequest struct {
	TenantID domain.TenantID `json:"-" uri:"tenant_id" binding:"required"`
	ID       int64           `uri:"id" binding:"required"`
	Hard     bool            `form:"hard,default=false"`
}

type ClearRecycleBinRequest struct {
	TenantID domain.TenantID `json:"-" uri:"tenant_id" binding:"required"`
	ID       int64           `uri:"id" binding:"required"`
}

type RestoreFromRecycleBinRequest struct {
	TenantID domain.TenantID `json:"-" uri:"tenant_id" binding:"required"`
	ID       int64           `uri:"id" binding:"required"`
}

type ListRequest struct {
	TenantID   domain.TenantID `json:"-" uri:"tenant_id" binding:"required"`
	KeyWord    string          `form:"keyword" binding:"max=20"`
	CategoryID int64           `form:"category_id"`
	Deleted    bool            `form:"deleted,default=false"`
	PageSize   int             `form:"page_size,default=5" binding:"min=5,max=50"`
	Page       int             `form:"page,default=1" binding:"min=1"`
}

type ImgListResponse struct {
	Total int64          `json:"total"`
	List  []*ImgResponse `json:"list"`
}

type CategoryResponse struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Prefix    string `json:"prefix"`
	CreatedAt int64  `json:"created_at"`
}

type CreateCategoryRequest struct {
	TenantID domain.TenantID `json:"-" uri:"tenant_id" binding:"required"`
	Title    string          `json:"title" binding:"required,max=10"`
	Prefix   string          `json:"prefix" binding:"required,max=20,slug"`
}

type UpdateCategoryRequest struct {
	ID       int64           `uri:"id" binding:"required"`
	TenantID domain.TenantID `json:"-" uri:"tenant_id" binding:"required"`
	Title    string          `json:"title" binding:"max=10"`
	Prefix   string          `json:"prefix" binding:"max=20"`
}

type DeleteCategoryRequest struct {
	ID       int64           `uri:"id" binding:"required"`
	TenantID domain.TenantID `json:"-" uri:"tenant_id" binding:"required"`
}

type ListCategoryRequest struct {
	TenantID domain.TenantID `json:"-" uri:"tenant_id" binding:"required"`
}
