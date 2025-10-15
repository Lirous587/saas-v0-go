package service

import (
	"saas/internal/comment/domain"
	"saas/internal/comment/templates"
)

const newCommentSubject = "新的评论信息"
const needAuditSubject = "新的评论需要审核"
const auditPassSubject = "评论审核通过"
const auditRejectSubject = "评论审核未通过"

func (s *service) sentCommentEmail(commentUser *domain.UserInfo, to string, relatedURL string, content string) error {
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

func (s *service) sentAuditPassEmail(to string, relatedURL string, content string) error {
	data := struct {
		CommentContent string
		RelatedURL     string
	}{
		CommentContent: content,
		RelatedURL:     relatedURL,
	}

	return s.mailer.SendWithTemplate(
		to,
		auditPassSubject,
		templates.TemplateAuditPass,
		data,
	)
}

func (s *service) sentAuditRejectEmail(to string, relatedURL string, content string) error {
	data := struct {
		CommentContent string
		RelatedURL     string
	}{
		CommentContent: content,
		RelatedURL:     relatedURL,
	}

	return s.mailer.SendWithTemplate(
		to,
		auditRejectSubject,
		templates.TemplateAuditReject,
		data,
	)
}
