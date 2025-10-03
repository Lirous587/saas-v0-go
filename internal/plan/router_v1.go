package plan

import (
	"github.com/gin-gonic/gin"
	"saas/internal/common/middleware/auth"
	"saas/internal/plan/handler"
)

func RegisterV1(r *gin.RouterGroup, handler *handler.HttpHandler) func() {
	g := r.Group("/v1/plan")
	{
		g.GET("", handler.List)
	}

	protect := g.Use(auth.JWTValidate(), auth.CasbinValited())
	{
		protect.POST("", handler.Create)
		protect.DELETE("/:id", handler.Delete)
		protect.PUT("/:id", handler.Update)
	}
	return nil
}
