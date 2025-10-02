//go:build wireinject
// +build wireinject

package role

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"saas/internal/role/adapters"
	"saas/internal/role/handler"
	"saas/internal/role/service"
)

func InitV1(r *gin.RouterGroup) func() {
	wire.Build(
		RegisterV1,
		handler.NewHttpHandler,
		service.NewRoleService,
		adapters.NewPSQLRoleRepository,
	)

	return nil
}
