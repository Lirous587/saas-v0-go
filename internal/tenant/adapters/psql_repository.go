package adapters

import (
	"context"
	"database/sql"
	"fmt"
	"saas/internal/common/orm"
	"saas/internal/common/reskit/codes"
	"saas/internal/tenant/domain"
	"slices"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/pkg/errors"
)

type TenantPSQLRepository struct {
}

func NewTenantPSQLRepository() domain.TenantRepository {
	return &TenantPSQLRepository{}
}

func (repo *TenantPSQLRepository) BeginTx(option ...*sql.TxOptions) (*sql.Tx, error) {
	var op *sql.TxOptions
	if len(option) > 1 {
		op = option[0]
	} else {
		op = &sql.TxOptions{
			Isolation: 0,
			ReadOnly:  false,
		}
	}
	return boil.BeginTx(context.TODO(), op)
}

func (repo *TenantPSQLRepository) InsertTx(tx *sql.Tx, tenant *domain.Tenant) (*domain.Tenant, error) {
	ormTenant := domainTenantToORM(tenant)

	if err := ormTenant.Insert(tx, boil.Infer()); err != nil {
		return nil, err
	}

	return ormTenantToDomain(ormTenant), nil
}

func (repo *TenantPSQLRepository) Update(tenant *domain.Tenant) error {
	ormTenant := domainTenantToORM(tenant)

	rows, err := ormTenant.UpdateG(boil.Infer())

	if err != nil {
		return err
	}
	if rows == 0 {
		return codes.ErrTenantNotFound
	}

	return nil
}

func (repo *TenantPSQLRepository) Delete(id int64) error {
	ormTenant := orm.Tenant{
		ID: id,
	}
	rows, err := ormTenant.DeleteG()

	if err != nil {
		return err
	}
	if rows == 0 {
		return codes.ErrTenantNotFound
	}
	return nil
}

func (repo *TenantPSQLRepository) Paging(query *domain.TenantPagingQuery) (*domain.TenantPagination, error) {
	mods := make([]qm.QueryMod, 0, 5)

	// 基本条件
	mods = append(mods, orm.TenantWhere.CreatorID.EQ(query.CreatorID))

	fetchLimit := query.PageSize + 1

	if query.Keyword != "" {
		like := "%" + query.Keyword + "%"
		mods = append(mods, qm.Where(fmt.Sprintf("(%s LIKE ? OR %s LIKE ?)", orm.TenantColumns.Name, orm.TenantColumns.Description), like, like))
	}

	// 游标与排序
	if query.AfterID > 0 {
		mods = append(mods, orm.TenantWhere.ID.GT(query.AfterID))
		mods = append(mods, qm.OrderBy(orm.TenantColumns.ID+" ASC"))
	} else if query.BeforeID > 0 {
		mods = append(mods, orm.TenantWhere.ID.LT(query.BeforeID))
		// 为了取到"上一页"的正确数据，先按 DESC 取最新的 N 条，然后反转切片返回给调用方（保持 ASC 展示）
		mods = append(mods, qm.OrderBy(orm.TenantColumns.ID+" DESC"))
	} else {
		mods = append(mods, qm.OrderBy(orm.TenantColumns.ID+" ASC"))
	}

	// limit
	mods = append(mods, qm.Limit(fetchLimit))

	ormTenants, err := orm.Tenants(mods...).AllG()
	if err != nil {
		return nil, err
	}

	hasMore := len(ormTenants) > query.PageSize

	// 截取
	if hasMore {
		ormTenants = ormTenants[:query.PageSize]
	}

	// 如果是 Before 分页，先按 DESC 取，结果需要反转为 ASC 展示
	if query.BeforeID > 0 && len(ormTenants) > 0 {
		slices.Reverse(ormTenants)
	}

	// 转换为 domain
	items := ormTenantsToDomain(ormTenants)

	// 基于请求方向和是否多取一条来推断 hasPrev/hasNext，避免额外 DB 查询
	isAfter := query.AfterID > 0
	isBefore := query.BeforeID > 0

	var hasPrev, hasNext bool
	switch {
	case isAfter:
		// 请求为 after（向后翻页），代表存在上一页（客户端传了游标）
		hasPrev = true
		hasNext = hasMore
	case isBefore:
		// 请求为 before（向前翻页），代表存在下一页（客户端传了游标）
		hasPrev = hasMore
		hasNext = true
	default:
		// 首页
		hasPrev = false
		hasNext = hasMore
	}

	return &domain.TenantPagination{
		Items:   items,
		HasNext: hasNext,
		HasPrev: hasPrev,
	}, nil
}

func (repo *TenantPSQLRepository) ExistSameName(creatorID int64, name string) (bool, error) {
	exist, err := orm.Tenants(
		orm.TenantWhere.CreatorID.EQ(creatorID),
		orm.TenantWhere.Name.EQ(name),
	).ExistsG()

	if err != nil {
		return false, errors.WithStack(err)
	}

	return exist, nil
}
