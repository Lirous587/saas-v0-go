//go:build wireinject
// +build wireinject

package comment

import (
	"saas/internal/comment/adapters"
	"saas/internal/comment/handler"
	"saas/internal/comment/service"
	"saas/internal/comment/templates"
	"saas/internal/common/email"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitV1(r *gin.RouterGroup) func() {
	wire.Build(
		RegisterV1,
		handler.NewHttpHandler,
		service.NewCommentService,
		adapters.NewCommentPSQLRepository,
		adapters.NewCommentRedisCache,
		email.NewMailer,
		templates.LoadCommentTemplates,
	)

	return nil
}
