package comment

import (
    "saas/internal/common/middleware/auth"
    "saas/internal/comment/handler"
	"github.com/gin-gonic/gin"
)

func RegisterV1(r *gin.RouterGroup, handler *handler.HttpHandler) func() {
	g := r.Group("/v1/comment")
	{
		g.GET("/:id",handler.Read)
		g.GET("", handler.List)
	}

    protect := g.Use(auth.JWTValidate())
    {
        protect.POST("", handler.Create)
        protect.DELETE("/:id", handler.Delete)
        protect.PUT("/:id", handler.Update)
    }
	return nil
}
