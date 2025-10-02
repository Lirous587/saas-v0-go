package adapters

import (
	"database/sql"
	"fmt"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/pkg/errors"
	"saas/internal/common/orm"
	"saas/internal/common/reskit/codes"
	"saas/internal/common/utils"
	"saas/internal/role/domain"
)

type RolePSQLRepository struct {
}

func NewRolePSQLRepository() domain.RoleRepository {
	return &RolePSQLRepository{}
}

func (repo *RolePSQLRepository) FindByID(id int64) (*domain.Role, error) {
	ormRole, err := orm.FindRoleG(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, codes.ErrRoleNotFound
		}
		return nil, err
	}
	return ormRoleToDomain(ormRole), nil
}

func (repo *RolePSQLRepository) Create(role *domain.Role) (*domain.Role, error) {
	ormRole := domainRoleToORM(role)

	if err := ormRole.InsertG(boil.Infer()); err != nil {
		return nil, err
	}

	return ormRoleToDomain(ormRole), nil
}

func (repo *RolePSQLRepository) Update(role *domain.Role) (*domain.Role, error) {
	ormRole := domainRoleToORM(role)

	rows, err := ormRole.UpdateG(boil.Infer())

	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, codes.ErrRoleNotFound
	}

	return ormRoleToDomain(ormRole), nil
}

func (repo *RolePSQLRepository) Delete(id int64) error {
	ormRole := orm.Role{
		ID: id,
	}
	rows, err := ormRole.DeleteG()

	if err != nil {
		return err
	}
	if rows == 0 {
		return codes.ErrRoleNotFound
	}
	return nil
}

func (repo *RolePSQLRepository) List(query *domain.RoleQuery) (*domain.RoleList, error) {
	var whereMods []qm.QueryMod
	if query.Keyword != "" {
		like := "%" + query.Keyword + "%"
		whereMods = append(whereMods, qm.Where(fmt.Sprintf("(%s LIKE ? OR %s LIKE ?)", orm.RoleColumns.Name, orm.RoleColumns.Description), like, like))
	}
	// 1.计算total
	total, err := orm.Roles(whereMods...).CountG()
	if err != nil {
		return nil, err
	}

	// 2.计算offset
	offset, err := utils.ComputeOffset(query.Page, query.PageSize)
	if err != nil {
		return nil, err
	}

	listMods := append(whereMods, qm.Offset(offset), qm.Limit(query.PageSize))

	// 3.查询数据
	role, err := orm.Roles(listMods...).AllG()
	if err != nil {
		return nil, err
	}

	return &domain.RoleList{
		Total: total,
		List:  ormRolesToDomain(role),
	}, nil
}

func (repo *RolePSQLRepository) FindUserRoleInTenant(userID, tenantID int64) (*domain.Role, error) {
	ut, err := orm.UserTenants(
		qm.Where(fmt.Sprintf("%s = ? AND %s = ?", orm.UserTenantColumns.UserID, orm.UserTenantColumns.TenantID), userID, tenantID),
		qm.Load(orm.UserTenantRels.Role),
	).OneG()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, codes.ErrRoleNotFound
		}
		return nil, err
	}

	// 检查 Role 是否为空
	if ut.R == nil || ut.R.Role == nil {
		return nil, codes.ErrRoleNotFound
	}

	return ormRoleToDomain(ut.R.Role), nil
}
