package service

import (
	"saas/internal/comment/domain"
	"saas/internal/comment/templates"
	"time"
)

const newCommentSubject = "新的评论信息"
const needAuditSubject = "评论审核通知"

func (s *service) sentCommentEmail(commentUser *domain.UserInfo, to string, relatedURL string, content string) error {

	data := struct {
		CommentUserNickname string
		CommentContent      string
		RelatedURL          string
		CreatedAt           int64
	}{
		CommentUserNickname: commentUser.NickName,
		CommentContent:      content,
		RelatedURL:          relatedURL,
		CreatedAt:           time.Now().Unix(),
	}

	return s.mailer.SendWithTemplate(
		to,
		newCommentSubject,
		templates.TemplateComment,
		data,
	)
}

func (s *service) sentNeedAuditEmail(commentUser *domain.UserInfo, to string, relatedURL string, content string) error {
	data := struct {
		CommentUserNickname string
		CommentContent      string
		RelatedURL          string
	}{
		CommentUserNickname: commentUser.NickName,
		CommentContent:      content,
		RelatedURL:          relatedURL,
	}

	return s.mailer.SendWithTemplate(
		to,
		needAuditSubject,
		templates.TemplateNeedAudit,
		data,
	)
}
