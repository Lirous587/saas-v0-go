package adapters

import (
	"context"
	"database/sql"
	"fmt"
	"saas/internal/comment/domain"
	"saas/internal/common/orm"
	"saas/internal/common/reskit/codes"
	"saas/internal/common/utils"
	"saas/internal/common/utils/dbkit"

	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/pkg/errors"
)

type CommentPSQLRepository struct {
}

func NewCommentPSQLRepository() domain.CommentRepository {
	return &CommentPSQLRepository{}
}

func (repo *CommentPSQLRepository) GetByID(tenantID domain.TenantID, commentID domain.CommentID) (*domain.Comment, error) {
	ormComment, err := orm.Comments(
		orm.CommentWhere.TenantID.EQ(string(tenantID)),
		orm.CommentWhere.ID.EQ(string(commentID)),
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

func (repo *CommentPSQLRepository) Delete(tenantID domain.TenantID, commentID domain.CommentID) error {
	rows, err := orm.Comments(
		orm.CommentWhere.TenantID.EQ(string(tenantID)),
		orm.CommentWhere.ID.EQ(string(commentID)),
	).DeleteAllG()

	if err != nil {
		return err
	}

	if rows == 0 {
		return codes.ErrCommentNotFound
	}

	return nil
}

func (repo *CommentPSQLRepository) Approve(tenantID domain.TenantID, commentID domain.CommentID) error {
	ormComment := &orm.Comment{
		TenantID: string(tenantID),
		ID:       string(commentID),
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

type replyCount struct {
	RootID string `boil:"root_id"`
	Count  int64  `boil:"reply_count"`
}

const (
	replyCountSelect = "COUNT(*) AS reply_count"
)

func (repo *CommentPSQLRepository) ListRoots(query *domain.CommentRootsQuery) ([]*domain.CommentRoot, error) {
	mods := make([]qm.QueryMod, 0, 9)
	mods = append(mods, orm.CommentWhere.TenantID.EQ(string(query.TenantID)))
	mods = append(mods, orm.CommentWhere.PlateID.EQ(string(query.PlateID)))
	// 评论状态为approved
	mods = append(mods, orm.CommentWhere.Status.EQ(orm.CommentStatusApproved))
	// 根评论过滤：root_id is null parent_id is null
	mods = append(mods, orm.CommentWhere.RootID.IsNull())
	mods = append(mods, orm.CommentWhere.ParentID.IsNull())

	// 使用游标代替offest
	if query.LastID != "" {
		mods = append(mods, orm.CommentWhere.ID.GT(string(query.LastID)))
	}

	// order 按ID升序排序
	mods = append(mods, qm.OrderBy(orm.CommentColumns.ID+" ASC"))

	// limit
	mods = append(mods, qm.Limit(query.PageSize))

	// 连表加载用户
	mods = append(mods, qm.Load(orm.CommentRels.User))

	// 执行查询
	ormComments, err := orm.Comments(mods...).AllG()
	if err != nil {
		return nil, err
	}

	// 为所有根评论计算回复数
	rootIDs := make([]string, 0, len(ormComments))
	for _, c := range ormComments {
		rootIDs = append(rootIDs, c.ID)
	}

	repliesCountMap := make(map[string]int64)
	if len(rootIDs) > 0 {
		var replyCounts []replyCount
		err := orm.NewQuery(
			qm.Select(orm.CommentTableColumns.RootID, replyCountSelect),
			qm.From(orm.TableNames.Comments),
			orm.CommentWhere.TenantID.EQ(string(query.TenantID)),
			orm.CommentWhere.PlateID.EQ(string(query.PlateID)),
			orm.CommentWhere.Status.EQ(orm.CommentStatusApproved),
			orm.CommentWhere.RootID.IN(rootIDs),
			qm.GroupBy(orm.CommentTableColumns.RootID),
		).BindG(context.Background(), &replyCounts)

		if err != nil {
			return nil, errors.WithStack(err)
		}

		for _, r := range replyCounts {
			repliesCountMap[r.RootID] = r.Count
		}
	}

	// 转换结果
	roots := make([]*domain.CommentRoot, 0, len(ormComments))
	for _, ormComment := range ormComments {
		userInfo := ormUserToDomain(ormComment.R.User)
		commentWithUser := &domain.CommentWithUser{
			ID:        domain.CommentID(ormComment.ID),
			User:      userInfo,
			ParentID:  domain.CommentID(ormComment.ParentID.String),
			RootID:    domain.CommentID(ormComment.RootID.String),
			Content:   ormComment.Content,
			LikeCount: ormComment.LikeCount,
			CreatedAt: ormComment.CreatedAt,
		}
		roots = append(roots, &domain.CommentRoot{
			CommentWithUser: commentWithUser,
			RepliesCount:    repliesCountMap[ormComment.ID],
		})
	}

	return roots, nil
}

func (repo *CommentPSQLRepository) ListReplies(query *domain.CommentRepliesQuery) ([]*domain.CommentReply, error) {
	mods := make([]qm.QueryMod, 0, 8)
	mods = append(mods, orm.CommentWhere.TenantID.EQ(string(query.TenantID)))
	mods = append(mods, orm.CommentWhere.PlateID.EQ(string(query.PlateID)))
	// 评论状态为approved
	mods = append(mods, orm.CommentWhere.Status.EQ(orm.CommentStatusApproved))
	// 根root条件
	mods = append(mods, orm.CommentWhere.RootID.EQ(null.StringFrom(string(query.RootID))))

	// 使用游标代替offest
	if query.LastID != "" {
		mods = append(mods, orm.CommentWhere.ID.GT(string(query.LastID)))
	}

	// order 按ID升序排序（确保游标有效）
	mods = append(mods, qm.OrderBy(orm.CommentColumns.ID+" ASC"))

	// limit
	mods = append(mods, qm.Limit(query.PageSize))

	// 连表加载用户
	mods = append(mods, qm.Load(orm.CommentRels.User))

	// 执行查询
	ormComments, err := orm.Comments(mods...).AllG()
	if err != nil {
		return nil, err
	}

	// 转换结果
	replies := make([]*domain.CommentReply, 0, len(ormComments))
	for i := range ormComments {
		userInfo := ormUserToDomain(ormComments[i].R.User)

		commentWithUser := &domain.CommentWithUser{
			ID:        domain.CommentID(ormComments[i].ID),
			User:      userInfo,
			ParentID:  domain.CommentID(ormComments[i].ParentID.String),
			RootID:    domain.CommentID(ormComments[i].RootID.String),
			Content:   ormComments[i].Content,
			LikeCount: ormComments[i].LikeCount,
			CreatedAt: ormComments[i].CreatedAt,
			// IsLiked:   false, // 待补充
		}

		replies = append(replies, &domain.CommentReply{
			CommentWithUser: commentWithUser,
		})
	}

	return replies, nil
}

func (repo *CommentPSQLRepository) UpdateLikeCount(tenantID domain.TenantID, commentID domain.CommentID, isLike bool) error {
	tx, err := boil.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	deltaStr := "+ 1"
	if !isLike {
		deltaStr = "- 1"
	}

	likeResStr := orm.CommentColumns.LikeCount + deltaStr

	sql := fmt.Sprintf(
		"UPDATE %s SET %s = %s WHERE id = $1 AND tenant_id = $2",
		orm.TableNames.Comments,
		orm.CommentColumns.LikeCount,
		likeResStr)

	result, err := queries.Raw(sql, commentID, tenantID).Exec(tx)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return codes.ErrCommentNotFound
	}

	return tx.Commit()
}

func (repo *CommentPSQLRepository) GetCommentUser(tenantID domain.TenantID, commentID domain.CommentID) (domain.UserID, error) {
	ormComment, err := orm.Comments(
		orm.CommentWhere.TenantID.EQ(string(tenantID)),
		orm.CommentWhere.ID.EQ(string(commentID)),
		qm.Select(orm.CommentColumns.UserID),
	).OneG()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", codes.ErrCommentNotFound
		}

		return "", err
	}

	return domain.UserID(ormComment.UserID), nil
}

func (repo *CommentPSQLRepository) GetUserIDsByRootORParent(tenantID domain.TenantID, plateID domain.PlateID, rootID domain.CommentID, parentID domain.CommentID) ([]domain.UserID, error) {
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

	userIDs := make([]domain.UserID, 0, 2)
	for i := range comments {
		userIDs = append(userIDs, domain.UserID(comments[i].UserID))
	}

	return userIDs, nil
}

func (repo *CommentPSQLRepository) GetTenantCreator(tenantID domain.TenantID) (*domain.UserInfo, error) {
	tenantUser, err := orm.Tenants(
		orm.TenantWhere.ID.EQ(string(tenantID)),
		qm.Select(orm.TenantColumns.CreatorID),
		qm.Load(orm.TenantRels.Creator),
	).OneG()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, codes.ErrTenantNotFound
		}
		return nil, err
	}

	// 获取关联的用户
	user := tenantUser.R.Creator
	if user == nil {
		return nil, errors.WithStack(codes.ErrUserNotFound.WithSlug("关联用户不存在"))
	}

	// 填充 UserInfo
	userInfo := &domain.UserInfo{
		ID:       domain.UserID(user.ID),
		NickName: user.Nickname,
	}

	userInfo.SetEmail(user.Email)

	return userInfo, nil
}

func (repo *CommentPSQLRepository) GetUserInfosByIDs(userIDs []domain.UserID) ([]*domain.UserInfo, error) {
	stringIDs := domain.UserIDs(userIDs).ToStringSlice()

	ormUsers, err := orm.Users(
		qm.WhereIn(fmt.Sprintf("%s in ?", orm.UserColumns.ID), utils.StringSliceToInterface(stringIDs)...),
		qm.Select(orm.UserColumns.ID, orm.UserColumns.Nickname, orm.UserColumns.Email),
	).AllG()

	if err != nil {
		return nil, err
	}

	return ormUsersToDomain(ormUsers), nil
}

func (repo *CommentPSQLRepository) GetUserInfoByID(userID domain.UserID) (*domain.UserInfo, error) {
	ormUser, err := orm.Users(
		orm.UserWhere.ID.EQ(string(userID)),
		qm.Select(orm.UserColumns.ID, orm.UserColumns.Nickname, orm.UserColumns.AvatarURL, orm.UserColumns.Email),
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
		boil.Infer(),
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

func (repo *CommentPSQLRepository) DeletePlate(tenantID domain.TenantID, plateID domain.PlateID) error {
	rows, err := orm.CommentPlates(
		orm.CommentPlateWhere.TenantID.EQ(string(tenantID)),
		orm.CommentPlateWhere.ID.EQ(string(plateID)),
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
	offset, err := dbkit.ComputeOffset(query.Page, query.PageSize)
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

func (repo *CommentPSQLRepository) GetPlateBelongByID(plateID domain.PlateID) (*domain.PlateBelong, error) {
	plate, err := orm.CommentPlates(
		orm.CommentPlateWhere.ID.EQ(string(plateID)),
		qm.Select(orm.CommentPlateColumns.ID, orm.CommentPlateColumns.BelongKey),
	).OneG()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, codes.ErrCommentPlateNotFound
		}
		return nil, err
	}

	return &domain.PlateBelong{
		ID:        domain.PlateID(plate.ID),
		BelongKey: plate.BelongKey,
	}, nil
}

func (repo *CommentPSQLRepository) GetPlateBelongByKey(tenantID domain.TenantID, belongKey string) (*domain.PlateBelong, error) {
	plate, err := orm.CommentPlates(
		orm.CommentPlateWhere.TenantID.EQ(string(tenantID)),
		orm.CommentPlateWhere.BelongKey.EQ(belongKey),
		qm.Select(orm.CommentPlateColumns.ID, orm.CommentPlateColumns.BelongKey),
	).OneG()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, codes.ErrCommentPlateNotFound
		}
		return nil, err
	}

	return &domain.PlateBelong{
		ID:        domain.PlateID(plate.ID),
		BelongKey: plate.BelongKey,
	}, nil
}

func (repo *CommentPSQLRepository) GetPlateRelatedURlByID(tenantID domain.TenantID, plateID domain.PlateID) (string, error) {
	plate, err := orm.CommentPlates(
		orm.CommentPlateWhere.TenantID.EQ(string(tenantID)),
		orm.CommentPlateWhere.ID.EQ(string(plateID)),
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

func (repo *CommentPSQLRepository) GetPlateConfig(tenantID domain.TenantID, plateID domain.PlateID) (*domain.PlateConfig, error) {
	ormConfig, err := orm.CommentPlateConfigs(
		orm.CommentPlateConfigWhere.TenantID.EQ(string(tenantID)),
		orm.CommentPlateConfigWhere.PlateID.EQ(string(plateID)),
	).OneG()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, codes.ErrCommentPlateConfigNotFound
		}
		return nil, err
	}

	return ormPlateConfigToDomain(ormConfig), nil
}
