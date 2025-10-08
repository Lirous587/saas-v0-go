package server

import (
	"saas/internal/common/reskit/codes"
	"strconv"

	"github.com/gin-gonic/gin"
)

const UserIDKey = "user_id"

func GetUserID(ctx *gin.Context) (int64, error) {
	uidStr, exist := ctx.Get(UserIDKey)
	if !exist {
		return 0, codes.ErrUserNotFound
	}

	userID, ok := uidStr.(int64)
	if !ok {
		return 0, codes.ErrUserIDInvalid
	}

	return userID, nil
}

const TenantIDKey = "tenant_id"

func SetTenantID(key string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从请求头获取
		if tenantIDStr := ctx.GetHeader(key); tenantIDStr != "" {
			if tenantID, err := strconv.ParseInt(tenantIDStr, 10, 64); err == nil {
				ctx.Set(TenantIDKey, tenantID)
				return
			}
		}

		// 从路径参数获取
		if tenantIDStr := ctx.Param(key); tenantIDStr != "" {
			if tenantID, err := strconv.ParseInt(tenantIDStr, 10, 64); err == nil {
				ctx.Set(TenantIDKey, tenantID)
				return
			}
		}

		// 从查询参数获取
		if tenantIDStr := ctx.Query(key); tenantIDStr != "" {
			if tenantID, err := strconv.ParseInt(tenantIDStr, 10, 64); err == nil {
				ctx.Set(TenantIDKey, tenantID)
				return
			}
		}

	}
}

func GetTenantID(ctx *gin.Context) (int64, error) {
	if tenantIDStr := ctx.Param(TenantIDKey); tenantIDStr != "" {
		tenantID, err := strconv.ParseInt(tenantIDStr, 10, 64)
		if err != nil {
			return 0, codes.ErrInvalidRequest.WithCause(err)
		}
		return tenantID, nil
	}

	if tid, exists := ctx.Get(TenantIDKey); exists {
		if id, ok := tid.(int64); ok {
			return id, nil
		}
	}

	return 0, codes.ErrInvalidRequest.WithSlug("无租户id来源")
}
