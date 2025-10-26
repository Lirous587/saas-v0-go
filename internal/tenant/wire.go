//go:build wireinject
// +build wireinject

package tenant

import (
	"saas/internal/tenant/adapters"
	"saas/internal/tenant/handler"
	"saas/internal/tenant/service"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitV1(r *gin.RouterGroup) func() {
	wire.Build(
		RegisterV1,
		handler.NewHttpHandler,
		service.NewTenantService,
		adapters.NewTenantPSQLRepository,
		adapters.NewTenantRedisCache,
	)

	return nil
}
