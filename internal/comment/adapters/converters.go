package adapters

import (
	"saas/internal/comment/domain"
	"saas/internal/common/orm"
)

func domainCommentToORM(comment *domain.Comment) *orm.Comment {
	if comment == nil {
		return nil
	}

	// 非null项
	ormComment := &orm.Comment{
		ID: comment.ID,
		// Title:     		comment.Title,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}

	// 处理null项
	// if comment.Description != "" {
	//  	ormComment.Description = null.StringFrom(comment.Description)
	// 	ormComment.Description.Valid = true
	// }

	return ormComment
}

func ormCommentToDomain(ormComment *orm.Comment) *domain.Comment {
	if ormComment == nil {
		return nil
	}

	// 非null项
	comment := &domain.Comment{
		ID: ormComment.ID,
		// Title:     		ormComment.Title,
		CreatedAt: ormComment.CreatedAt,
		UpdatedAt: ormComment.UpdatedAt,
	}

	// 处理null项
	// if ormComment.Description.Valid {
	//  	comment.Description = ormComment.Description.String
	// }

	return comment
}

func ormCommentsToDomain(ormComments []*orm.Comment) []*domain.Comment {
	if len(ormComments) == 0 {
		return nil
	}

	comments := make([]*domain.Comment, 0, len(ormComments))
	for _, ormComment := range ormComments {
		if ormComment != nil {
			comments = append(comments, ormCommentToDomain(ormComment))
		}
	}
	return comments
}
