package tenant

import (
	"saas/internal/common/middleware/auth"
	"saas/internal/common/server"
	"saas/internal/tenant/handler"

	"github.com/gin-gonic/gin"
)

func RegisterV1(r *gin.RouterGroup, handler *handler.HttpHandler) func() {
	g := r.Group("/v1/tenant")

	protect := g.Use(auth.JWTValidate())
	{
		// todo 创建租户 目前未接入交易中间件
		protect.POST("", handler.Create)
		protect.GET("", handler.ListByKeyset)
		protect.GET("/check_name", handler.CheckName)
	}

	// 仅租户创建者可访问的路由
	creatorOnly := g.Group("", auth.JWTValidate(), server.SetTenantID("id"), auth.TenantCreatorValited())
	{
		creatorOnly.GET("/:id", handler.Read)
		creatorOnly.PUT("/:id", handler.Update)
		creatorOnly.GET("/:id/plan", handler.GetPlan)
	}
	return nil
}
