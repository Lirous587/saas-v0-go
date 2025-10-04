//go:build wireinject
// +build wireinject

package tenant

import (
	"saas/internal/common/email"
	planAdapters "saas/internal/plan/adapters"
	planService "saas/internal/plan/service"
	roleAdapters "saas/internal/role/adapters"
	roleService "saas/internal/role/service"
	"saas/internal/tenant/adapters"
	"saas/internal/tenant/handler"
	"saas/internal/tenant/service"
	"saas/internal/tenant/templates"
	userAdapters "saas/internal/user/adapters"

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

		email.NewMailer,
		templates.LoadTenantTemplates,

		// plan服务
		planAdapters.NewPlanPSQLRepository,
		planService.NewPlanService,

		// role服务
		roleAdapters.NewRolePSQLRepository,
		roleAdapters.NewRoleRedisCache,
		roleService.NewRoleService,

		// user repo
		userAdapters.NewUserPSQLRepository,
	)

	return nil
}
