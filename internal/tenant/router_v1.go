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
		// todo 目前未接入交易中间件
		protect.POST("", handler.Create)
		protect.GET("/:id", server.SetTenantID("id"), auth.CasbinValited(), handler.Read)
		protect.POST("/:id/gen_invite_token", server.SetTenantID("id"), auth.CasbinValited(), handler.GenInviteToken)
		//protect.DELETE("/:id", handler.Delete)
		//protect.PUT("/:id", handler.Update)
	}
	return nil
}
