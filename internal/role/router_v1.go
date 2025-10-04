package role

import (
	"saas/internal/common/middleware/auth"
	"saas/internal/role/handler"

	"github.com/gin-gonic/gin"
)

func RegisterV1(r *gin.RouterGroup, handler *handler.HttpHandler) func() {
	g := r.Group("/v1/role")

	protect := g.Use(auth.JWTValidate(), auth.CasbinValited())
	{
		protect.GET("/:tenant_id", handler.List)

		// todo 系统管理员方可编辑角色
		// protect.POST("/:tenant_id/:id", handler.Create)
		// protect.DELETE("/:tenant_id/:id", handler.Delete)
		// protect.PUT("/:tenant_id/:id", handler.Update)
	}
	return nil
}
