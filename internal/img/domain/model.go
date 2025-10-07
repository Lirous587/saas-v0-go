package domain

import (
	"time"
)

type TenantID int64

type Img struct {
	ID           int64
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

type ImgQuery struct {
	TenantID   TenantID
	CategoryID int64
	Keyword    string
	Page       int
	PageSize   int
	Deleted    bool
}

type ImgList struct {
	List  []*Img
	Total int64
}

type Category struct {
	ID        int64
	TenantID  TenantID
	Title     string
	Prefix    string
	CreatedAt time.Time
}

type R2Config struct {
	TenantID        int64
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
