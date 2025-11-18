package dbkit

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/aarondl/sqlboiler/v4/queries/qm"
)

/*
Keyset 使用示例（Feed，带 score 排序）

示例 1 — 简单 Feed（created_at 为主游标，score 为显示优先级但不参与 existence 检查）
1) 创建 keyset：
   ks := dbkit.NewKeyset[*domain.FeedItem, time.Time](
       orm.CommentColumns.ID,
       orm.CommentColumns.CreatedAt,
       query.PrevCursor,
       query.NextCursor,
       query.PageSize,
       dbkit.WithPrimaryOrder[*domain.FeedItem, time.Time](dbkit.SortDesc),
       dbkit.WithExtraOrderCols[*domain.FeedItem, time.Time]("score DESC"), // 仅用于排序
   )

2) 构造基础 filters（与 existence 保持一致）：
   baseMods := []qm.QueryMod{
       orm.CommentWhere.Status.EQ(orm.CommentStatusApproved),
       // ... 其他过滤 ...
   }

3) 加入 keyset mods 并查询：
   mods := ks.ApplyKeysetMods(baseMods)
   ormItems, _ := orm.Comments(mods...).AllG()
   items := ormFeedItemsFromORM(ormItems)

4) 精确 existence（仅用 primary+id）：
   // 若 extraOrderCols 只是展示 tie-breaker，且 primary+id 已能区分记录，
   // 则可直接用 BeforeWhere/AfterWhere（基于 primary+id）
   exists := func(primary time.Time, id string, checkPrev bool) (bool, error) {
       var cond qm.QueryMod
       if checkPrev {
           cond = ks.BeforeWhere(primary, id)
       } else {
           cond = ks.AfterWhere(primary, id)
       }
       checkMods := append(baseMods, cond, qm.Limit(1))
       return orm.Comments(checkMods...).ExistsG()
   }

5) 构造分页结果：
   res, _ := ks.BuildPaginationResultWithExistence(items, exists)


示例 2 — 复杂 Feed（created_at + score 共同排序，需精确 hasPrev/hasNext）
说明：如果你把 score 放进排序列里，并期望精确判断“首/尾之外是否有记录”，则 existence 查询必须也用相同的元组比较 (created_at, score, id)。
两种实现方法：
A) 读取首/尾记录的 score（从 item 直接取，或通过 DB 再取），并在 exists 中用三元组比较：
   exists := func(primary time.Time, id string, checkPrev bool) (bool, error) {
       // 从数据库查出所需的 score（或从 repo 的 domain item 获得）
       score, err := repo.GetScoreByID(id)
       if err != nil { return false, err }

       var cond qm.QueryMod
       if checkPrev {
           cond = qm.Where("(created_at, score, id) > (?, ?, ?)", primary, score, id)
       } else {
           cond = qm.Where("(created_at, score, id) < (?, ?, ?)", primary, score, id)
       }
       checkMods := append(baseMods, cond, qm.Limit(1))
       return orm.Comments(checkMods...).ExistsG()
   }
注意：此实现会产生额外 DB 查询来获得 score（一次/端点存在性判断），代价较小但需要注意性能。

B) 更优雅但更改量更大：把 score 一并包含到游标（扩展 keysetCursor 结构与 GetCursorPrimary 返回类型 P），使 keyset 的 encode/decode 带上 score，这样 BeforeWhere/AfterWhere 可以直接构造 "(created_at, score, id) < ... / > ..." 而无需额外查询。实现需要修改 keyset 以支持 composite primary 值（例如 struct{CreatedAt time.Time; Score float64}）并把 ApplyKeysetMods、BeforeWhere/AfterWhere 调整为支持多列比较。

总结：
- 推荐初期用示例 A（简单且易实现）。
- 如果对性能敏感且需要稳定无额外查询的判断，考虑实现 B（把额外排序列包含进游标）。
- 无论哪种方式，务必保证：ApplyKeysetMods 的排序列、Before/AfterWhere 的比较逻辑与 existence 的实现保持一致，否则 hasPrev/hasNext 将不准确。
*/

type SortDirection string

const (
	SortAsc  SortDirection = "ASC"
	SortDesc SortDirection = "DESC"
)

type KeysetOption[T CursorFields[P], P any] func(k *keyset[T, P])

func WithPrimaryOrder[T CursorFields[P], P any](dir SortDirection) KeysetOption[T, P] {
	return func(k *keyset[T, P]) { k.primaryOrder = dir }
}
func WithIDOrder[T CursorFields[P], P any](dir SortDirection) KeysetOption[T, P] {
	return func(k *keyset[T, P]) { k.idOrder = dir }
}

type keysetCursor[P any] struct {
	Primary P      `json:"primary"`
	ID      string `json:"id"`
}

type CursorFields[P any] interface {
	GetCursorPrimary() P
	GetID() string
}

type keyset[T CursorFields[P], P any] struct {
	IDCol        string
	idOrder      SortDirection
	PrimaryCol   string
	primaryOrder SortDirection

	PrevCursor string
	NextCursor string
	PageSize   int
}

type paginationResult[T any] struct {
	Items      []T
	PrevCursor string
	NextCursor string
	HasPrev    bool
	HasNext    bool
}

func NewKeyset[T CursorFields[P], P any](
	idCol string,
	primaryCol string,
	prevCursor string,
	nextCursor string,
	pageSize int,
	opts ...KeysetOption[T, P],
) *keyset[T, P] {
	if pageSize <= 0 {
		pageSize = 5
	}
	k := &keyset[T, P]{
		IDCol:        idCol,
		PrimaryCol:   primaryCol,
		PrevCursor:   prevCursor,
		NextCursor:   nextCursor,
		PageSize:     pageSize,
		primaryOrder: SortDesc, // feed 常见默认降序按时间
		// idOrder 默认不设置，随后与 primaryOrder 同步
		idOrder: "",
	}
	for _, opt := range opts {
		opt(k)
	}

	// 若未显式设置 idOrder，则默认与 primaryOrder 保持一致
	if string(k.idOrder) == "" {
		k.idOrder = k.primaryOrder
	}
	return k
}

func (k keyset[T, P]) encode(cursor *keysetCursor[P]) string {
	if cursor == nil {
		return ""
	}
	data, _ := json.Marshal(cursor)
	return base64.StdEncoding.EncodeToString(data)
}

func (k keyset[T, P]) decodeCursor(cursorStr string) *keysetCursor[P] {
	if cursorStr == "" {
		return nil
	}
	data, err := base64.StdEncoding.DecodeString(cursorStr)
	if err != nil {
		return nil
	}
	var payload keysetCursor[P]
	if json.Unmarshal(data, &payload) != nil {
		return nil
	}
	return &payload
}

func (k keyset[T, P]) extractKeysetCursor(item T) *keysetCursor[P] {
	return &keysetCursor[P]{
		Primary: item.GetCursorPrimary(),
		ID:      item.GetID(),
	}
}

// ApplyKeysetMods 附加keyset查询的mods
func (k keyset[T, P]) ApplyKeysetMods(base []qm.QueryMod) []qm.QueryMod {
	var cursor *keysetCursor[P]
	var orderParts []string

	// ORDER BY 初始配置
	pOrder := string(k.primaryOrder)
	idOrder := string(k.idOrder)

	if k.NextCursor != "" {
		cursor = k.decodeCursor(k.NextCursor)
		if cursor != nil {
			// 使用 AfterWhere 保证和 After/Before 的方向一致
			base = append(base, k.AfterWhere(cursor.Primary, cursor.ID))
		}
		// next 保持 order 与配置一致
		pOrder = string(k.primaryOrder)
		idOrder = string(k.idOrder)
	} else if k.PrevCursor != "" {
		cursor = k.decodeCursor(k.PrevCursor)
		if cursor != nil {
			// 使用 BeforeWhere 保证和 After/Before 的方向一致
			base = append(base, k.BeforeWhere(cursor.Primary, cursor.ID))
		}
		// Prev 查询时主列、ID 的排序方向都要反转
		if k.primaryOrder == SortAsc {
			pOrder = string(SortDesc)
		} else {
			pOrder = string(SortAsc)
		}
		if k.idOrder == SortAsc {
			idOrder = string(SortDesc)
		} else {
			idOrder = string(SortAsc)
		}
	} else {
		// 无游标，按配置方向
		pOrder = string(k.primaryOrder)
		idOrder = string(k.idOrder)
	}

	orderParts = append(orderParts, fmt.Sprintf("%s %s", k.PrimaryCol, pOrder))
	orderParts = append(orderParts, fmt.Sprintf("%s %s", k.IDCol, idOrder))
	order := strings.Join(orderParts, ", ")

	base = append(base, qm.OrderBy(order))
	base = append(base, qm.Limit(k.PageSize+1))

	return base
}

// BeforeWhere 返回用于判断是否存在“在给定 (primary,id) 之前（prev page）”的 qm.Where
// 语义："之前" 是相对于 ApplyKeysetMods 所使用的主排序方向（primaryOrder）。
func (k keyset[T, P]) BeforeWhere(primary P, id string) qm.QueryMod {
	// Before = 获取“在游标之前”的项目（prev page）
	// 若主列 DESC，则“之前”是更新的项 => > ，否则 <。
	if k.primaryOrder == SortDesc {
		return qm.Where(fmt.Sprintf("(%s, %s) > (?, ?)", k.PrimaryCol, k.IDCol), primary, id)
	}
	return qm.Where(fmt.Sprintf("(%s, %s) < (?, ?)", k.PrimaryCol, k.IDCol), primary, id)
}

// AfterWhere 返回用于判断是否存在“在给定 (primary,id) 之后（next page）”的 qm.Where
func (k keyset[T, P]) AfterWhere(primary P, id string) qm.QueryMod {
	// After = 获取“在游标之后”的项目（next page）
	// 若主列 DESC，则“之后”是更旧的项 => < ，否则 >。
	if k.primaryOrder == SortDesc {
		return qm.Where(fmt.Sprintf("(%s, %s) < (?, ?)", k.PrimaryCol, k.IDCol), primary, id)
	}
	return qm.Where(fmt.Sprintf("(%s, %s) > (?, ?)", k.PrimaryCol, k.IDCol), primary, id)
}

// BuildWithExistence 通过用户提供的 existence 检查器做精确的 hasPrev/hasNext 判断。
// exists 期望实现：func(primary P, id string, checkPrev bool) (bool, error)
// - checkPrev=true 检查是否存在比 (primary,id) 更“前”的记录（previous page）
// - checkPrev=false 检查是否存在比 (primary,id) 更“后”的记录（next page）
func (k keyset[T, P]) BuildWithExistence(domainSlice []T, exists func(P, string, bool) (bool, error)) (*paginationResult[T], error) {
	hasMore := len(domainSlice) > k.PageSize

	// 截取
	if hasMore {
		domainSlice = domainSlice[:k.PageSize]
	}

	isPrev := k.PrevCursor != ""

	// 如果是 Prev 分页，先按 DESC 取，结果需要反转为 ASC 展示
	if isPrev && len(domainSlice) > 0 {
		slices.Reverse(domainSlice)
	}

	// 生成新的游标
	var prevCursor, nextCursor string
	if len(domainSlice) > 0 {
		first := domainSlice[0]
		last := domainSlice[len(domainSlice)-1]

		prevCursor = k.encode(k.extractKeysetCursor(first))
		nextCursor = k.encode(k.extractKeysetCursor(last))
	}

	// 精确判断 hasPrev/hasNext：交给 exists 回调
	var hasPrev, hasNext bool
	if len(domainSlice) == 0 {
		hasPrev = false
		hasNext = false
	} else {
		first := domainSlice[0]
		last := domainSlice[len(domainSlice)-1]

		// 由 caller 提供实现，需与 ApplyKeysetMods 里的比较方向一致
		pFirst := k.extractKeysetCursor(first).Primary
		idFirst := k.extractKeysetCursor(first).ID
		pLast := k.extractKeysetCursor(last).Primary
		idLast := k.extractKeysetCursor(last).ID

		var err error
		hasPrev, err = exists(pFirst, idFirst, true)
		if err != nil {
			return nil, err
		}
		hasNext, err = exists(pLast, idLast, false)
		if err != nil {
			return nil, err
		}
	}

	// 没有下一页时 next_cursor置空
	if !hasNext {
		nextCursor = ""
	}
	// 没有上一页时 prev_cursor置空
	if !hasPrev {
		prevCursor = ""
	}

	return &paginationResult[T]{
		Items:      domainSlice,
		PrevCursor: prevCursor,
		NextCursor: nextCursor,
		HasPrev:    hasPrev,
		HasNext:    hasNext,
	}, nil
}

// BuildFeed: 轻量 feed-only 策略，专注 next-page（不做存在性查询）
// - 总是返回 PrevCursor（方便前端尝试上一页）
// - HasPrev 采用简单推断：当请求带有 next_cursor 时，认为上一页存在（用户刚刚翻到下一页）
// - HasNext 基于 pageSize+1 判断；Prev 分页时仍会反转结果以展示正确顺序
func (k keyset[T, P]) BuildFeed(domainSlice []T) *paginationResult[T] {
	hasMore := len(domainSlice) > k.PageSize
	if hasMore {
		domainSlice = domainSlice[:k.PageSize]
	}

	// 如果是 Prev 分页（极少使用），反转回正常显示顺序
	isPrev := k.PrevCursor != ""
	if isPrev && len(domainSlice) > 0 {
		slices.Reverse(domainSlice)
	}

	var prevCursor, nextCursor string
	if len(domainSlice) > 0 {
		first := domainSlice[0]
		last := domainSlice[len(domainSlice)-1]
		// Always build prev cursor so frontend can try to go back
		prevCursor = k.encode(k.extractKeysetCursor(first))
		// next cursor only present when we actually have more items
		if hasMore {
			nextCursor = k.encode(k.extractKeysetCursor(last))
		}
	}

	// Feed-only 语义：
	// - HasPrev: 如果客户端带有 next_cursor（即刚点“下一页”），认为上页还可以回溯
	// - HasNext: 依据 pageSize+1 判定（精确）
	hasPrev := k.NextCursor != ""
	hasNext := hasMore

	return &paginationResult[T]{
		Items:      domainSlice,
		PrevCursor: prevCursor,
		NextCursor: nextCursor,
		HasPrev:    hasPrev,
		HasNext:    hasNext,
	}
}
