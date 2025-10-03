//go:build wireinject
// +build wireinject

package tenant

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	planAdapters "saas/internal/plan/adapters"
	planService "saas/internal/plan/service"
	roleAdapters "saas/internal/role/adapters"
	roleService "saas/internal/role/service"
	"saas/internal/tenant/adapters"
	"saas/internal/tenant/handler"
	"saas/internal/tenant/service"
)

func InitV1(r *gin.RouterGroup) func() {
	wire.Build(
		RegisterV1,
		handler.NewHttpHandler,
		service.NewTenantService,
		adapters.NewTenantPSQLRepository,

		// plan服务
		planAdapters.NewPlanPSQLRepository,
		planService.NewPlanService,

		// role服务
		roleAdapters.NewRolePSQLRepository,
		roleAdapters.NewRoleRedisCache,
		roleService.NewRoleService,
	)

	return nil
}
