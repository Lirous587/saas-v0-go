package dbkit

import (
	"fmt"
	"slices"
	"strings"

	"github.com/aarondl/sqlboiler/v4/queries/qm"
)

// SortField 定义排序字段和方向
type SortField struct {
	Column    string // 数据库列名
	Direction string // "ASC" 或 "DESC"
}

// KeysetCursor 表示分页游标值
type KeysetCursor = []interface{}

type CursorEncoder[T any] interface {
	Encode(res *T) string
}

type CursorDecoder interface {
	DecodePrev() KeysetCursor
	DecodeNext() KeysetCursor
}

// Keyset 支持多字段排序的分页结构
type Keyset[T any] struct {
	BeforeCursor KeysetCursor
	AfterCursor  KeysetCursor
	PageSize     int
	SortFields   []SortField
	encoder      CursorEncoder[T]
}

func NewKeyset[T any](pageSize int, sortFields []SortField, decoder CursorDecoder, encoder CursorEncoder[T]) *Keyset[T] {
	if pageSize <= 0 {
		pageSize = 5
	}
	return &Keyset[T]{
		PageSize:     pageSize,
		SortFields:   sortFields,
		BeforeCursor: decoder.DecodePrev(),
		AfterCursor:  decoder.DecodeNext(),
		encoder:      encoder,
	}
}

func (k *Keyset[T]) ApplyKeysetMods(base []qm.QueryMod) []qm.QueryMod {
	if k.AfterCursor != nil {
		base = append(base, k.buildWhereCondition(false))
	} else if k.BeforeCursor != nil {
		base = append(base, k.buildWhereCondition(true))
	}
	orderByClause := k.buildOrderByClause()
	base = append(base, qm.OrderBy(orderByClause))

	fetchLimit := k.PageSize + 1
	base = append(base, qm.Limit(fetchLimit))
	return base
}

func (k *Keyset[T]) buildWhereCondition(isBefore bool) qm.QueryMod {
	cursor := k.AfterCursor
	if isBefore && k.BeforeCursor != nil {
		cursor = k.BeforeCursor
	}

	clauses := make([]string, 0, len(k.SortFields))
	values := make([]interface{}, 0)

	for i := range k.SortFields {
		parts := make([]string, 0, i+1)
		for j := 0; j < i; j++ {
			parts = append(parts, fmt.Sprintf("%s = ?", k.SortFields[j].Column))
			values = append(values, cursor[j])
		}
		parts = append(parts, fmt.Sprintf("%s %s ?", k.SortFields[i].Column, k.comparator(i, isBefore)))
		values = append(values, cursor[i])
		clauses = append(clauses, "("+strings.Join(parts, " AND ")+")")
	}

	return qm.Where(strings.Join(clauses, " OR "), values...)
}

func (k *Keyset[T]) comparator(idx int, isBefore bool) string {
	dir := strings.ToUpper(k.SortFields[idx].Direction)
	op := ">"
	if dir == "DESC" {
		op = "<"
	}
	if isBefore {
		if op == ">" {
			op = "<"
		} else {
			op = ">"
		}
	}
	return op
}

// buildOrderByClause 构建ORDER BY子句
func (k *Keyset[T]) buildOrderByClause() string {
	clauses := make([]string, len(k.SortFields))
	for i, field := range k.SortFields {
		// 如果是上一页请求，反转排序方向
		direction := field.Direction
		if k.BeforeCursor != nil {
			if direction == "ASC" {
				direction = "DESC"
			} else {
				direction = "ASC"
			}
		}
		clauses[i] = fmt.Sprintf("%s %s", field.Column, direction)
	}
	return strings.Join(clauses, ", ")
}

type PaginationResult[T any] struct {
	Items      []*T
	HasNext    bool
	HasPrev    bool
	NextCursor string
	PrevCursor string
}

// BuildPaginationResult 构建分页结果
func (k *Keyset[T]) BuildPaginationResult(items []*T) *PaginationResult[T] {
	hasMore := len(items) > k.PageSize

	// 截取结果
	if hasMore {
		items = items[:k.PageSize]
	}

	// 如果是上一页请求，反转结果
	if k.BeforeCursor != nil && len(items) > 0 {
		slices.Reverse(items)
	}

	// 计算游标
	var nextCursor, prevCursor string
	if len(items) > 0 {
		firstItem := items[0]
		lastItem := items[len(items)-1]

		prevCursor = k.encoder.Encode(firstItem)
		nextCursor = k.encoder.Encode(lastItem)
	}

	// 判断分页状态
	isAfter := k.AfterCursor != nil
	isBefore := k.BeforeCursor != nil

	var hasPrev, hasNext bool
	switch {
	case isAfter:
		hasPrev = true
		hasNext = hasMore
	case isBefore:
		hasPrev = hasMore
		hasNext = true
	default:
		hasPrev = false
		hasNext = hasMore
	}

	return &PaginationResult[T]{
		Items:      items,
		HasPrev:    hasPrev,
		HasNext:    hasNext,
		NextCursor: nextCursor,
		PrevCursor: prevCursor,
	}
}
