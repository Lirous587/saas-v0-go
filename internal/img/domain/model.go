package domain

import (
	"time"
)

type Img struct {
	ID           ImgID
	TenantID     TenantID
	Path         string
	Description  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time
	publicPreURL string
}

func (img *Img) IsDeleted() bool {
	// deleted_time非零值则为软删除记录
	return !img.DeletedAt.IsZero()
}

// CanDeleted 状态转换 当前图片是否可以软删除或硬删除
func (img *Img) CanDeleted() bool {
	return !img.IsDeleted()
}

func (img *Img) SetPublicPreURL(preURL string) {
	img.publicPreURL = preURL
}

func (img *Img) GetPublicPreURL() string {
	return img.publicPreURL
}

func (img *Img) GetID() string {
	return img.ID.String()
}

func (img *Img) GetCursorPrimary() time.Time {
	return img.CreatedAt
}

type ListByKeysetQuery struct {
	TenantID   TenantID
	CategoryID CategoryID
	PrevCursor string
	NextCursor string
	Keyword    string
	PageSize   int
	Deleted    bool
}

type ListByKeysetResult struct {
	Items      []*Img
	PrevCursor string
	NextCursor string
	HasPrev    bool
	HasNext    bool
}

type Category struct {
	ID        CategoryID
	TenantID  TenantID
	Title     string
	Prefix    string
	CreatedAt time.Time
}

type R2Config struct {
	TenantID        TenantID
	AccountID       string
	AccessKeyID     string
	secretAccessKey string
	PublicBucket    string
	PublicURLPrefix string
	DeleteBucket    string
}

func (r *R2Config) GetSecretAccessKey() string {
	return r.secretAccessKey
}

func (r *R2Config) SetSecretAccessKey(key string) {
	r.secretAccessKey = key
}

type R2SecretAccessKey string
