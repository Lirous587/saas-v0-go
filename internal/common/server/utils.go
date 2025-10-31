package server

import (
	"saas/internal/common/reskit/codes"

	"github.com/gin-gonic/gin"
)

const UserIDKey = "user_id"

func GetUserID(ctx *gin.Context) (string, error) {
	uidStr, exist := ctx.Get(UserIDKey)
	if !exist {
		return "", codes.ErrUserNotFound
	}

	userID, ok := uidStr.(string)
	if !ok {
		return "", codes.ErrUserIDInvalid
	}

	return userID, nil
}

const TenantIDKey = "tenant_id"

func SetTenantID(key string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从请求头获取
		if tenantID := ctx.GetHeader(key); tenantID != "" {
			ctx.Set(TenantIDKey, tenantID)
		}

		// 从路径参数获取
		if tenantID := ctx.Param(key); tenantID != "" {
			ctx.Set(TenantIDKey, tenantID)
		}

		// 从查询参数获取
		if tenantID := ctx.Query(key); tenantID != "" {
			ctx.Set(TenantIDKey, tenantID)
		}

	}
}

func GetTenantID(ctx *gin.Context) (string, error) {
	if tenantID := ctx.Param(TenantIDKey); tenantID != "" {
		return tenantID, nil
	}

	if tid, exists := ctx.Get(TenantIDKey); exists {
		if id, ok := tid.(string); ok {
			return id, nil
		}
	}

	return "", codes.ErrInvalidRequest.WithSlug("无租户id来源")
}
