package utils

import (
	"time"

	"github.com/aarondl/sqlboiler/v4/queries/qm"
)

// Cursor 表示复合游标 (created_at, id)
type Cursor struct {
	CreatedAt *time.Time
	ID        int64
}

// BuildCursorMods 根据传入的 after/before 游标与列名，构建查询模块（where/order/limit）
// 返回的 needReverse 表示当为 before 模式时，查询结果需在应用层反转为升序返回。
// fetchLimit = pageSize + 1，用于判断是否有更多页（调用方负责按 pageSize 截断）。
func BuildCursorMods(createdAtCol, idCol string, after, before *Cursor, pageSize int) (mods []qm.QueryMod, needReverse bool, fetchLimit int) {

	return
}
