package adapters

import (
	"database/sql"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/pkg/errors"
	"saas/internal/comment/domain"
	"saas/internal/common/orm"
	"saas/internal/common/reskit/codes"
	"saas/internal/common/utils"
)

type CommentPSQLRepository struct {
}

func NewCommentPSQLRepository() domain.CommentRepository {
	return &CommentPSQLRepository{}
}

func (repo *CommentPSQLRepository) FindByID(id int64) (*domain.Comment, error) {
	ormComment, err := orm.FindCommentG(id)
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

func (repo *CommentPSQLRepository) Update(comment *domain.Comment) (*domain.Comment, error) {
	ormComment := domainCommentToORM(comment)

	rows, err := ormComment.UpdateG(boil.Infer())

	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, codes.ErrCommentNotFound
	}

	return ormCommentToDomain(ormComment), nil
}

func (repo *CommentPSQLRepository) Delete(id int64) error {
	ormComment := orm.Comment{
		ID: id,
	}
	rows, err := ormComment.DeleteG()

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
	if query.Keyword != "" {
		// like := "%" + query.Keyword + "%"
		// whereMods = append(whereMods, qm.Where(fmt.Sprintf("(%s LIKE ? OR %s LIKE ?)", orm.CommentColumns.Title, orm.CommentColumns.Description), like, like))
	}
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
