package tenant

import (
	"saas/internal/common/middleware/auth"
	"saas/internal/tenant/handler"

	"github.com/gin-gonic/gin"
)

func RegisterV1(r *gin.RouterGroup, handler *handler.HttpHandler) func() {
	g := r.Group("/v1/tenant")

	protect := g.Use(auth.JWTValidate())
	{
		// todo 创建租户 目前未接入交易中间件
		protect.POST("", handler.Create)
		protect.PUT("/:id", handler.Update)
		protect.GET("", handler.Paging)
		protect.GET("/check_name", handler.CheckName)
	}
	return nil
}
