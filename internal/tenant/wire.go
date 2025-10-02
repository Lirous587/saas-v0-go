//go:build wireinject
// +build wireinject

package tenant

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"saas/internal/tenant/adapters"
	"saas/internal/tenant/handler"
	"saas/internal/tenant/service"
)

func InitV1(r *gin.RouterGroup) func() {
	wire.Build(
		RegisterV1,
		handler.NewHttpHandler,
		service.NewTenantService,
		adapters.NewPSQLTenantRepository,
	)

	return nil
}
