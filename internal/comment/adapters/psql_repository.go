package adapters

import (
	"database/sql"
	"fmt"
	"saas/internal/comment/domain"
	"saas/internal/common/orm"
	"saas/internal/common/reskit/codes"
	"saas/internal/common/utils"
	roleDomain "saas/internal/role/domain"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/pkg/errors"
)

type CommentPSQLRepository struct {
	role roleDomain.Role
}

func NewCommentPSQLRepository() domain.CommentRepository {
	return &CommentPSQLRepository{}
}

func (repo *CommentPSQLRepository) GetByID(tenantID domain.TenantID, id int64) (*domain.Comment, error) {
	ormComment, err := orm.Comments(
		qm.Where(fmt.Sprintf("%s = ?", orm.CommentColumns.TenantID), tenantID),
		qm.And(fmt.Sprintf("%s = ?", orm.CommentColumns.ID), id),
	).OneG()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, codes.ErrCommentNotFound
		}
		return nil, err
	}

	return ormCommentToDomain(ormComment), nil
}

func (repo *CommentPSQLRepository) Create(comment *domain.Comment) (*domain.Comment, error) {
	ormComment := domainCommentToORM(comment)

	if err := ormComment.InsertG(boil.Infer()); err != nil {
		return nil, err
	}

	return ormCommentToDomain(ormComment), nil
}

func (repo *CommentPSQLRepository) Delete(tenantID domain.TenantID, id int64) error {
	rows, err := orm.Comments(
		qm.Where(fmt.Sprintf("%s = ?", orm.CommentColumns.TenantID), tenantID),
		qm.And(fmt.Sprintf("%s = ?", orm.CommentColumns.ID), id),
	).DeleteAllG()

	if err != nil {
		return err
	}

	if rows == 0 {
		return codes.ErrCommentNotFound
	}

	return nil
}

func (repo *CommentPSQLRepository) Approve(tenantID domain.TenantID, id int64) error {
	ormComment := &orm.Comment{
		TenantID: int64(tenantID),
		ID:       id,
		Status:   orm.CommentStatusApproved,
	}

	rows, err := ormComment.UpdateG(boil.Whitelist(orm.CommentColumns.Status))
	if err != nil {
		return err
	}
	if rows == 0 {
		return codes.ErrCommentNotFound
	}

	return nil
}

func (repo *CommentPSQLRepository) List(query *domain.CommentQuery) (*domain.CommentList, error) {
	var whereMods []qm.QueryMod
	// if query.Keyword != "" {
	// 	// like := "%" + query.Keyword + "%"
	// 	// whereMods = append(whereMods, qm.Where(fmt.Sprintf("(%s LIKE ? OR %s LIKE ?)", orm.CommentColumns.Title, orm.CommentColumns.Description), like, like))
	// }
	// 1.计算total
	total, err := orm.Comments(whereMods...).CountG()
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
	comment, err := orm.Comments(listMods...).AllG()
	if err != nil {
		return nil, err
	}

	return &domain.CommentList{
		Total: total,
		List:  ormCommentsToDomain(comment),
	}, nil
}

func (repo *CommentPSQLRepository) GetCommentUser(tenantID domain.TenantID, commentID int64) (int64, error) {
	ormComment, err := orm.Comments(
		qm.Where(fmt.Sprintf("%s = ?", orm.CommentColumns.TenantID), tenantID),
		qm.And(fmt.Sprintf("%s = ?", orm.CommentColumns.ID), commentID),
		qm.Select(orm.CommentColumns.UserID),
	).OneG()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, codes.ErrCommentNotFound
		}

		return 0, err
	}

	return ormComment.UserID, nil
}

func (repo *CommentPSQLRepository) GetUserIdsByRootORParent(tenantID domain.TenantID, plateID int64, rootID int64, parentID int64) ([]int64, error) {
	comments, err := orm.Comments(
		qm.Where(fmt.Sprintf("%s = ? AND %s = ? AND (%s = ? OR %s = ?)",
			orm.CommentColumns.TenantID,
			orm.CommentColumns.PlateID,
			orm.CommentColumns.ID,
			orm.CommentColumns.ID),
			tenantID, plateID, rootID, parentID),
		qm.Select(orm.CommentColumns.UserID),
	).AllG()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, codes.ErrCommentNotFound
		}
		return nil, err
	}

	userIds := make([]int64, 0, 2)
	for i := range comments {
		userIds = append(userIds, comments[i].UserID)
	}

	return userIds, nil
}

func (repo *CommentPSQLRepository) GetDomainAdminByTenant(tenantID domain.TenantID) (*domain.UserInfo, error) {
	tenantUserRole, err := orm.TenantUserRoles(
		qm.Where(fmt.Sprintf("%s = ?", orm.TenantUserRoleColumns.TenantID), tenantID),
		qm.Where(fmt.Sprintf("%s = ?", orm.TenantUserRoleColumns.RoleID), repo.role.GetTenantadmin().ID),
		qm.Select(orm.TenantUserRoleColumns.UserID),
		qm.Load(orm.TenantUserRoleRels.User),
	).OneG()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, codes.ErrTenantNotFound
		}
		return nil, err
	}

	// 获取关联的用户
	user := tenantUserRole.R.User
	if user == nil {
		return nil, errors.WithStack(codes.ErrUserNotFound.WithSlug("关联用户不存在"))
	}

	// 填充 UserInfo
	userInfo := &domain.UserInfo{
		ID:       user.ID,
		NickName: user.Nickname,
	}

	userInfo.SetEmail(user.Email)

	return userInfo, nil
}

func (repo *CommentPSQLRepository) GetUserInfosByIds(ids []int64) ([]*domain.UserInfo, error) {
	ormUsers, err := orm.Users(
		qm.WhereIn(fmt.Sprintf("%s in ?", orm.UserColumns.ID), utils.Int64SliceToInterface(ids)...),
		qm.Select(orm.UserColumns.ID, orm.UserColumns.Nickname, orm.UserColumns.Email),
	).AllG()

	if err != nil {
		return nil, err
	}

	return ormUsersToDomain(ormUsers), nil
}

func (repo *CommentPSQLRepository) GetUserInfoByID(id int64) (*domain.UserInfo, error) {
	ormUser, err := orm.Users(
		qm.Where(fmt.Sprintf("%s = ?", orm.UserColumns.ID), id),
		qm.Select(orm.UserColumns.ID, orm.UserColumns.Nickname, orm.UserColumns.Avatar, orm.UserColumns.Email),
	).OneG()

	if err != nil {
		return nil, err
	}

	return ormUserToDomain(ormUser), nil
}

func (repo *CommentPSQLRepository) SetTenantConfig(config *domain.TenantConfig) error {
	ormConfig := domainTenantConfigToORM(config)
	if err := ormConfig.UpsertG(
		true,
		[]string{orm.CommentTenantConfigColumns.TenantID},
		boil.Blacklist(
			orm.CommentTenantConfigColumns.ClientToken,
		),
		boil.Greylist( //使用GreyList 因为用Infer的话IfAudit设置为false时不会生效 此时就会导致使用默认值true 与请求冲突
			orm.CommentTenantConfigColumns.IfAudit,
		),
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (repo *CommentPSQLRepository) GetTenantConfig(tenantID domain.TenantID) (*domain.TenantConfig, error) {
	ormConfig, err := orm.CommentTenantConfigs(
		qm.Where(fmt.Sprintf("%s = ?", orm.CommentTenantConfigColumns.TenantID), tenantID),
	).OneG()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, codes.ErrCommentTenantConfigNotFound
		}
		return nil, err
	}

	return ormTenantConfigToDomain(ormConfig), nil
}

func (repo *CommentPSQLRepository) ExistTenantConfigByID(tenantID domain.TenantID) (bool, error) {
	exist, err := orm.CommentTenantConfigs(
		qm.Where(fmt.Sprintf("%s = ?", orm.CommentTenantConfigColumns.TenantID), tenantID),
	).ExistsG()

	if err != nil {
		return false, err
	}

	return exist, nil
}

func (repo *CommentPSQLRepository) CreatePlate(plate *domain.Plate) error {
	ormPlate := domainPlateToORM(plate)

	return ormPlate.InsertG(boil.Infer())
}

func (repo *CommentPSQLRepository) UpdatePlate(plate *domain.Plate) error {
	ormPlate := domainPlateToORM(plate)

	rows, err := ormPlate.UpdateG(boil.Infer())
	if err != nil {
		return err
	}

	if rows == 0 {
		return codes.ErrCommentPlateNotFound
	}

	return nil
}

func (repo *CommentPSQLRepository) DeletePlate(tenantID domain.TenantID, id int64) error {
	rows, err := orm.CommentPlates(
		qm.Where(fmt.Sprintf("%s = ?", orm.CommentPlateColumns.TenantID), tenantID),
		qm.Where(fmt.Sprintf("%s = ?", orm.CommentPlateColumns.ID), id),
	).DeleteAllG()

	if err != nil {
		return err
	}

	if rows == 0 {
		return codes.ErrCommentPlateNotFound
	}

	return nil
}

func (repo *CommentPSQLRepository) ListPlate(query *domain.PlateQuery) (*domain.PlateList, error) {
	var whereMods []qm.QueryMod
	if query.Keyword != "" {
		like := "%" + query.Keyword + "%"
		whereMods = append(whereMods, qm.Where(fmt.Sprintf("(%s LIKE ? OR %s LIKE ?)", orm.CommentPlateColumns.BelongKey, orm.CommentPlateColumns.Summary), like, like))
	}
	// 1.计算total
	total, err := orm.CommentPlates(whereMods...).CountG()
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
	plate, err := orm.CommentPlates(listMods...).AllG()
	if err != nil {
		return nil, err
	}

	return &domain.PlateList{
		Total: total,
		List:  ormPlatesToDomain(plate),
	}, nil
}

func (repo *CommentPSQLRepository) ExistPlateBykey(tenantID domain.TenantID, belongKey string) (bool, error) {
	exist, err := orm.CommentPlates(
		qm.Where(fmt.Sprintf("%s = ?", orm.CommentPlateColumns.TenantID), tenantID),
		qm.And(fmt.Sprintf("%s = ?", orm.CommentPlateColumns.BelongKey), belongKey),
	).ExistsG()

	if err != nil {
		return false, err
	}

	return exist, nil
}

func (repo *CommentPSQLRepository) GetPlateBelongByID(id int64) (*domain.PlateBelong, error) {
	plate, err := orm.CommentPlates(
		qm.Where(fmt.Sprintf("%s = ?", orm.CommentPlateColumns.ID), id),
		qm.Select(orm.CommentPlateColumns.ID, orm.CommentPlateColumns.BelongKey),
	).OneG()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, codes.ErrCommentPlateNotFound
		}
		return nil, err
	}

	return &domain.PlateBelong{
		ID:        plate.ID,
		BelongKey: plate.BelongKey,
	}, nil
}

func (repo *CommentPSQLRepository) GetPlateBelongByKey(tenantID domain.TenantID, belongKey string) (*domain.PlateBelong, error) {
	plate, err := orm.CommentPlates(
		qm.Where(fmt.Sprintf("%s = ?", orm.CommentPlateColumns.TenantID), tenantID),
		qm.Where(fmt.Sprintf("%s = ?", orm.CommentPlateColumns.BelongKey), belongKey),
		qm.Select(orm.CommentPlateColumns.ID, orm.CommentPlateColumns.BelongKey),
	).OneG()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, codes.ErrCommentPlateNotFound
		}
		return nil, err
	}

	return &domain.PlateBelong{
		ID:        plate.ID,
		BelongKey: plate.BelongKey,
	}, nil
}

func (repo *CommentPSQLRepository) GetPlateRelatedURlByID(tenantID domain.TenantID, id int64) (string, error) {
	plate, err := orm.CommentPlates(
		qm.Where(fmt.Sprintf("%s = ?", orm.CommentPlateColumns.TenantID), tenantID),
		qm.Where(fmt.Sprintf("%s = ?", orm.CommentPlateColumns.ID), id),
		qm.Select(orm.CommentPlateColumns.RelatedURL),
	).OneG()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", codes.ErrCommentPlateNotFound
		}
		return "", err
	}

	return plate.RelatedURL, nil
}

func (repo *CommentPSQLRepository) SetPlateConfig(config *domain.PlateConfig) error {
	ormConfig := domainPlateConfigToORM(config)
	if err := ormConfig.UpsertG(
		true,
		[]string{orm.CommentPlateConfigColumns.TenantID, orm.CommentPlateConfigColumns.PlateID}, // 冲突列：复合主键的两个字段
		boil.Greylist(
			orm.CommentPlateConfigColumns.IfAudit,
		),
		boil.Greylist(
			orm.CommentPlateConfigColumns.IfAudit,
		),
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (repo *CommentPSQLRepository) GetPlateConfig(tenantID domain.TenantID, plateID int64) (*domain.PlateConfig, error) {
	ormConfig, err := orm.CommentPlateConfigs(
		qm.Where(fmt.Sprintf("%s = ?", orm.CommentPlateConfigColumns.TenantID), tenantID),
		qm.Where(fmt.Sprintf("%s = ?", orm.CommentPlateConfigColumns.PlateID), plateID),
	).OneG()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, codes.ErrCommentPlateConfigNotFound
		}
		return nil, err
	}

	return ormPlateConfigToDomain(ormConfig), nil
}
