package comment

import (
	"github.com/gin-gonic/gin"
	"saas/internal/comment/handler"
	"saas/internal/common/middleware/auth"
)

func RegisterV1(r *gin.RouterGroup, handler *handler.HttpHandler) func() {
	// 0.应该要给每个租户分发令牌 避免刷接口导致的所有评论数据泄露 稍后实现

	// 查询权限 用不同路由去做
	// 访客仅仅允许分页查询
	// domain_admin允许高级查询

	// api
	// 评论
	// 删除评论
	// 审计 (用户可选是否开启) 审计和令牌管理和新建立一个表 以tenant_id benlong_key做隔离
	// 点赞评论 (低优先级)
	// 分页查询
	// 高级查询

	g := r.Group("/v1/comment/:tenant_id")
	{
		// 访客 获取评论
		// 获取根评论
		g.GET("/:belong_key/roots", auth.OptionalJWTValidate(), handler.ListRoots)
		// 根据根评论去获取其树下评论
		g.GET("/:belong_key/:root_id/replies", auth.OptionalJWTValidate(), handler.ListReplies)
	}

	protect := g.Group("", auth.JWTValidate())
	{
		// 创建评论
		protect.POST("/:belong_key", handler.Create)
		// 删除评论
		protect.DELETE("/:id", handler.Delete)

		// 低优先级 点赞/取消点赞
		protect.PUT("/like/:id", handler.ToggleLike)
	}

	// 仅租户创建者可访问的路由
	creatorOnly := g.Group("", auth.JWTValidate(), auth.TenantCreatorValited())
	{
		// 管理员
		// 审计
		creatorOnly.PUT("/:id", auth.TenantCreatorValited(), handler.Audit)

		// 全局配置
		creatorOnly.PUT("/config", handler.SetTenantConfig)
		creatorOnly.GET("/config", handler.GetTenantConfig)

		// 板块管理子组
		plateGroup := creatorOnly.Group("/plate")
		{
			plateGroup.POST("", handler.CreatePlate)
			plateGroup.PUT("/:id", handler.UpdatePlate)
			plateGroup.DELETE("/:id", handler.DeletePlate)
			plateGroup.GET("", handler.ListPlate)
			// 板块配置
			plateGroup.PUT("/config/:belong_key", handler.SetPlateConfig)
			plateGroup.GET("/config/:id", handler.GetPlateConfig)

			plateGroup.GET("/check_name", handler.CheckPlateBelongKey)
		}

	}

	return nil
}
