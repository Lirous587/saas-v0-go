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
		// 因为要使用url来实现 故此只能用get
		g.GET("/entry/:id", handler.Enter)
	}

	protect := g.Use(auth.JWTValidate())
	{
		// 此租户信息
		protect.GET("/:id", server.SetTenantID("id"), auth.CasbinValited(), handler.Read)
		// todo 创建租户 目前未接入交易中间件
		protect.POST("", handler.Create)
		// protect.PUT("/:id", handler.Update)
		protect.POST("/:id/gen_invite_token", server.SetTenantID("id"), auth.CasbinValited(), handler.GenInviteToken)
		protect.POST("/:id/invite", server.SetTenantID("id"), auth.CasbinValited(), handler.Invite)

		// 租户下的用户信息及其角色
		protect.GET("/:id/users", server.SetTenantID("id"), auth.CasbinValited(), handler.GetUsers)

		// todo 分配角色
		// protect.POST("/:id/:user_id", server.SetTenantID("id"), auth.CasbinValited(), handler.xx)
	}
	return nil
}
