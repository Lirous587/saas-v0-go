package tenant

import (
	"saas/internal/common/middleware/auth"
	"saas/internal/common/server"
	"saas/internal/tenant/handler"

	"github.com/gin-gonic/gin"
)

func RegisterV1(r *gin.RouterGroup, handler *handler.HttpHandler) func() {
	g := r.Group("/v1/tenant")
	{
		//g.GET("", handler.List)
	}

	protect := g.Use(auth.JWTValidate())
	{
		protect.POST("", handler.Create)
		protect.GET("/:id", server.SetTenantID("id"), auth.CasbinValited(), handler.Read)
		//protect.DELETE("/:id", handler.Delete)
		//protect.PUT("/:id", handler.Update)
	}
	return nil
}
