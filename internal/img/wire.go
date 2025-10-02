//go:build wireinject
// +build wireinject

package img

import (
	"saas/internal/img/adapters"
	"saas/internal/img/handler"
	"saas/internal/img/service"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitV1(r *gin.RouterGroup) func() {
	wire.Build(
		RegisterV1,
		handler.NewHttpHandler,
		service.NewImgService,
		adapters.NewImgPSQLRepository,
		adapters.NewImgRedisCache,
	)

	return nil
}
