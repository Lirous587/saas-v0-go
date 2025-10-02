//go:build wireinject
// +build wireinject

package user

import (
	"saas/internal/user/adapters"
	"saas/internal/user/handler"
	"saas/internal/user/service"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitV1(r *gin.RouterGroup) func() {
	wire.Build(
		RegisterV1,
		handler.NewHttpHandler,
		service.NewTokenService,
		service.NewUserService,
		adapters.NewUserPSQLRepository,
		adapters.NewTokenRedisCache,
	)
	return nil
}
