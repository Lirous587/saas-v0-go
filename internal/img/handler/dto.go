package handler

import "saas/internal/img/domain"

type ImgResponse struct {
	ID          domain.ImgID `json:"id"`
	URL         string       `json:"url"`
	Description string       `json:"description,omitempty"`
	CreatedAt   int64        `json:"created_at"`
	UpdatedAt   int64        `json:"updated_at"`
}

type UploadRequest struct {
	TenantID    domain.TenantID   `json:"-" uri:"tenant_id" binding:"required,uuid"`
	Path        string            `form:"path" binding:"omitempty,slug"`
	Description string            `form:"description" binding:"max=60"`
	CategoryID  domain.CategoryID `form:"category_id" binding:"omitempty,uuid"`
}

type DeleteRequest struct {
	TenantID domain.TenantID `json:"-" uri:"tenant_id" binding:"required,uuid"`
	ID       domain.ImgID    `uri:"id" binding:"required,uuid"`
	Hard     bool            `form:"hard,default=false"`
}

type ClearRecycleBinRequest struct {
	TenantID domain.TenantID `json:"-" uri:"tenant_id" binding:"required,uuid"`
	ID       domain.ImgID    `uri:"id" binding:"required,uuid"`
}

type RestoreFromRecycleBinRequest struct {
	TenantID domain.TenantID `json:"-" uri:"tenant_id" binding:"required,uuid"`
	ID       domain.ImgID    `uri:"id" binding:"required,uuid"`
}

type ListRequest struct {
	TenantID   domain.TenantID   `json:"-" uri:"tenant_id" binding:"required,uuid"`
	CategoryID domain.CategoryID `form:"category_id" binding:"omitempty,uuid"`
	KeyWord    string            `form:"keyword" binding:"max=20"`
	Deleted    bool              `form:"deleted,default=false"`
	PageSize   int               `form:"page_size,default=5" binding:"min=5,max=50"`
	Page       int               `form:"page,default=1" binding:"min=1"`
}

type ImgListResponse struct {
	Total int64          `json:"total"`
	List  []*ImgResponse `json:"list"`
}

type CategoryResponse struct {
	ID        domain.CategoryID `json:"id"`
	Title     string            `json:"title"`
	Prefix    string            `json:"prefix"`
	CreatedAt int64             `json:"created_at"`
}

type CreateCategoryRequest struct {
	TenantID domain.TenantID `json:"-" uri:"tenant_id" binding:"required,uuid"`
	Title    string          `json:"title" binding:"required,max=10"`
	Prefix   string          `json:"prefix" binding:"required,max=20,slug"`
}

type UpdateCategoryRequest struct {
	ID       domain.CategoryID `json:"-" uri:"id" binding:"required,uuid"`
	TenantID domain.TenantID   `json:"-" uri:"tenant_id" binding:"required,uuid"`
	Title    string            `json:"title" binding:"required,max=10"`
	Prefix   string            `json:"prefix" binding:"required,max=20,slug"`
}

type DeleteCategoryRequest struct {
	ID       domain.CategoryID `json:"-" uri:"id" binding:"required,uuid"`
	TenantID domain.TenantID   `json:"-" uri:"tenant_id" binding:"required,uuid"`
}

type ListCategoryRequest struct {
	TenantID domain.TenantID `json:"-" uri:"tenant_id" binding:"required,uuid"`
}

type SetR2ConfigRequest struct {
	TenantID        domain.TenantID `json:"-" uri:"tenant_id" binding:"required,uuid"`
	AccountID       string          `json:"account_id" binding:"required,len=32"`
	AccessKeyID     string          `json:"access_key_id" binding:"required,len=32"`
	PublicBucket    string          `json:"public_bucket" binding:"required,max=32"`
	PublicURLPrefix string          `json:"public_url_prefix" binding:"required,domain_url,max=128"`
	DeleteBucket    string          `json:"delete_bucket" binding:"required,max=32"`
}

type R2ConfigResponse struct {
	AccountID       string `json:"account_id"`
	AccessKeyID     string `json:"access_key_id"`
	PublicBucket    string `json:"public_bucket"`
	PublicURLPrefix string `json:"public_url_prefix"`
	DeleteBucket    string `json:"delete_bucket"`
}

type GetR2ConfigRequest struct {
	TenantID domain.TenantID `json:"-" uri:"tenant_id" binding:"required,uuid"`
}

type SetR2SecretAccessKeyRequest struct {
	TenantID        domain.TenantID          `json:"-" uri:"tenant_id" binding:"required,uuid"`
	SecretAccessKey domain.R2SecretAccessKey `json:"secret_access_key" binding:"required,len=64"`
}

type IsSetR2SecretRequest struct {
	TenantID domain.TenantID `json:"-" uri:"tenant_id" binding:"required,uuid"`
}

type IsSetR2SecretResponse struct {
	IsSet bool `json:"is_set"`
}
