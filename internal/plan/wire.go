//go:build wireinject
// +build wireinject

package plan

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"saas/internal/plan/adapters"
	"saas/internal/plan/handler"
	"saas/internal/plan/service"
)

func InitV1(r *gin.RouterGroup) func() {
	wire.Build(
		RegisterV1,
		handler.NewHttpHandler,
		service.NewPlanService,
		adapters.NewPlanPSQLRepository,
	)

	return nil
}
