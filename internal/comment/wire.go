//go:build wireinject
// +build wireinject

package comment

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"saas/internal/comment/adapters"
	"saas/internal/comment/handler"
	"saas/internal/comment/service"
)

func InitV1(r *gin.RouterGroup) func() {
	wire.Build(
		RegisterV1,
		handler.NewHttpHandler,
		service.NewCommentService,
		adapters.NewCommentPSQLRepository,
	)

	return nil
}
