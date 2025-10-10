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
		// 访客：分页查询
		g.GET("", handler.List)
	}

	protect := g.Use(auth.JWTValidate(), auth.CasbinValited())
	{
		// 用户：创建评论
		protect.POST("/:plate", handler.Create)

		// 用户：删除评论（只能删自己的，或管理员删任意）
		protect.DELETE("/:plate/:id", handler.Delete)

		// 低优先级：点赞/取消点赞
		// protect.POST("/:id/like", handler.Like)
		// protect.DELETE("/:id/like", handler.Unlike)

		// 管理员
		// 高级查询
		// protect.GET("/advanced", handler.AdvancedList)
		// 审计
		// protect.PUT("/:id")

		// 板块管理
		protect.POST("/plate", handler.CreatePlate)
		protect.DELETE("/plate/:id", handler.DeletePlate)
		protect.GET("/plate", handler.ListPlate)

		// 全局配置
		protect.PUT("/config", handler.SetTenantConfig)
		protect.GET("/config", handler.GetTenantConfig)
		// 板块配置
		protect.PUT("/:plate/config", handler.SetPlateConfig)
		protect.GET("/:plate/config", handler.GetPlateConfig)
	}

	return nil
}
