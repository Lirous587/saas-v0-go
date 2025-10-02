package role

import (
	"fmt"
	"saas/internal/common/middleware/auth"
	"saas/internal/common/server"
	"saas/internal/role/handler"

	"github.com/gin-gonic/gin"
)

func RegisterV1(r *gin.RouterGroup, handler *handler.HttpHandler) func() {
	g := r.Group("/v1/role")

	protect := g.Use(auth.JWTValidate())
	{
		protect.GET(fmt.Sprintf("/:%s", server.TenantIDKey), auth.CasbinValited(), handler.List)
		protect.POST("", handler.Create)
		protect.DELETE("/:id", handler.Delete)
		protect.PUT("/:id", handler.Update)
	}
	return nil
}
