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

	g := r.Group("/v1/comment")
	{
		// 访客：分页查询
		g.GET("", handler.List)
	}

	protect := g.Use(auth.JWTValidate(), auth.CasbinValited())
	{
		// 用户：创建评论
		protect.POST("/:tenant_id/:belong_key", handler.Create)

		// 用户：删除评论（只能删自己的，或管理员删任意）
		protect.DELETE("/:tenant_id/:belong_key/:id", handler.Delete)

		// 低优先级：点赞/取消点赞
		// protect.POST("/:id/like", handler.Like)
		// protect.DELETE("/:id/like", handler.Unlike)
	}

	// 管理员：高级查询（需额外权限检查，如 domain_admin）
	// admin := g.Use(auth.JWTValidate(), auth.CasbinValited())
	// {
	// admin.GET("/advanced", handler.AdvancedList) // 高级查询，如按状态、用户等过滤

	// 审计（用户可选开启）
	// admin.GET("/audit", handler.AuditList)       // 查询审计日志
	// admin.PUT("/audit/:id", handler.AuditUpdate) // 更新审核状态
	// }

	return nil
}
