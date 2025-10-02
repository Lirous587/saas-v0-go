package tenant

import (
	"github.com/gin-gonic/gin"
	"saas/internal/common/middleware/auth"
	"saas/internal/tenant/handler"
)

func RegisterV1(r *gin.RouterGroup, handler *handler.HttpHandler) func() {
	g := r.Group("/v1/tenant")
	{
		//g.GET("/:id", handler.Read)
		//g.GET("", handler.List)
	}

	protect := g.Use(auth.Validate())
	{
		protect.POST("", handler.Create)
		//protect.DELETE("/:id", handler.Delete)
		//protect.PUT("/:id", handler.Update)
	}
	return nil
}
