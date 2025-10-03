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

	protect := g.Use(auth.JWTValidate(), auth.CasbinValited())
	{
		protect.GET(fmt.Sprintf("/:%s", server.TenantIDKey), handler.List)
		protect.POST(fmt.Sprintf("/:%s", server.TenantIDKey), handler.Create)
		protect.DELETE(fmt.Sprintf("/:%s/:id", server.TenantIDKey), handler.Delete)
		protect.PUT(fmt.Sprintf("/:%s/:id", server.TenantIDKey), handler.Update)
	}
	return nil
}
