package domain

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

type Tenant struct {
	ID          int64
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CreatorID   int64
}

type TenantPagingQuery struct {
	PageSize     int
	CreatorID    int64
	AfterCursor  string
	BeforeCursor string
	Keyword      string
}

func (q *TenantPagingQuery) ToCursorDecoder() *TenantCursorDecoder {
	return &TenantCursorDecoder{
		BeforeRaw: q.BeforeCursor,
		AfterRaw:  q.AfterCursor,
	}
}

type TenantCursorEncoder struct{}

func (TenantCursorEncoder) Encode(t *Tenant) string {
	if t == nil {
		return ""
	}
	payload := []interface{}{t.UpdatedAt.UnixNano(), t.ID}
	data, _ := json.Marshal(payload)
	return base64.StdEncoding.EncodeToString(data)
}

type TenantCursorDecoder struct {
	BeforeRaw string
	AfterRaw  string
}

func (d *TenantCursorDecoder) DecodePrev() []interface{} {
	return decodeCursor(d.BeforeRaw)
}

func (d *TenantCursorDecoder) DecodeNext() []interface{} {
	return decodeCursor(d.AfterRaw)
}

func decodeCursor(raw string) []interface{} {
	if raw == "" {
		return nil
	}
	data, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return nil
	}
	var payload []interface{}
	if json.Unmarshal(data, &payload) != nil || len(payload) != 2 {
		return nil
	}
	ts, okTs := payload[0].(float64)
	id, okID := payload[1].(float64)
	if !okTs || !okID {
		return nil
	}
	return []interface{}{time.Unix(0, int64(ts)), int64(id)}
}

type TenantPagination struct {
	Items      []*Tenant
	PrevCursor string
	NextCursor string
	HasPrev    bool
	HasNext    bool
}
