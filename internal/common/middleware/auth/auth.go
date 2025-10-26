package auth

import (
	casbinadapter "saas/internal/common/casbin"
	"saas/internal/common/orm"
	"saas/internal/common/reskit/codes"
	"saas/internal/common/reskit/response"
	"saas/internal/common/server"

	useradapter "saas/internal/user/adapters"
	userdomain "saas/internal/user/domain"
	userService "saas/internal/user/service"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
)

var tokenServer userdomain.TokenService

var enforcer *casbin.Enforcer

func Init() {
	// 初始化token服务
	tokenCache := useradapter.NewTokenRedisCache()
	userRepo := useradapter.NewUserPSQLRepository()
	tokenServer = userService.NewTokenService(tokenCache, userRepo)

	// 初始化casbin服务
	var err error
	adapter := casbinadapter.NewSQLBoilerCasbinAdapter()
	enforcer, err = casbin.NewEnforcer("./model.conf", adapter)
	if err != nil {
		panic("创建执行器失败: " + err.Error())
	}

	err = enforcer.LoadPolicy()
	if err != nil {
		panic("加载策略失败:" + err.Error())
	}
}

const (
	authHeaderKey = "Authorization"
	bearerPrefix  = "Bearer "
)

// 解析 Authorization 头部的 Token
func parseTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader(authHeaderKey)
	if authHeader == "" {
		return "", errors.New("token为空")
	}

	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", errors.New("token格式错误")
	}

	return strings.TrimPrefix(authHeader, bearerPrefix), nil
}

func JWTValidate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从请求头解析 Token
		tokenStr, err := parseTokenFromHeader(c)
		if err != nil {
			response.Error(c, codes.ErrTokenFormatInvalid)
			return
		}

		// 2. 验证token
		isExpire, err := tokenServer.ValidateAccessToken(tokenStr)
		if err != nil {
			if isExpire {
				response.Error(c, codes.ErrTokenExpired)
			} else {
				response.Error(c, codes.ErrTokenInvalid)
			}
			return
		}

		// 3. 解析 Token
		payload, err := tokenServer.ParseAccessToken(tokenStr)
		if err != nil {
			response.Error(c, codes.ErrTokenInvalid)
			return
		}

		// 3. 将用户 相关信息存入上下文
		c.Set(server.UserIDKey, payload.UserID)

		c.Next()
	}
}

func OptionalJWTValidate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从请求头解析 Token
		tokenStr, err := parseTokenFromHeader(c)
		if err != nil {
			c.Next()
			return
		}

		// 2. 验证token
		_, err = tokenServer.ValidateAccessToken(tokenStr)
		if err != nil {
			c.Next()
			return
		}

		// 3. 解析 Token
		payload, err := tokenServer.ParseAccessToken(tokenStr)
		if err != nil {
			c.Next()
			return
		}

		// 3. 将用户 相关信息存入上下文
		c.Set(server.UserIDKey, payload.UserID)

		c.Next()
	}
}

func TenantCreatorValited() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取useID
		userID, err := server.GetUserID(ctx)
		if err != nil {
			response.Error(ctx, codes.ErrUnauthorized)
			ctx.Abort()
			return
		}

		// 获取tenantID
		tenantID, err := server.GetTenantID(ctx)
		if err != nil {
			response.Error(ctx, err)
			ctx.Abort()
			return
		}

		// 检测用户是否为该租户的创建者
		exist, err := orm.Tenants(
			orm.TenantWhere.ID.EQ(tenantID),
			orm.TenantWhere.CreatorID.EQ(userID),
		).ExistsG()
		if err != nil || !exist {
			response.Error(ctx, codes.ErrTenantNotCreator)
		}

		ctx.Next()
	}
}

const tenantAdmin = "domain_admin"

func CasbinValited() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 1.获取useID
		userID, err := server.GetUserID(ctx)
		if err != nil {
			response.Error(ctx, codes.ErrUnauthorized)
			ctx.Abort()
			return
		}

		// 2.获取tenantID
		tenantID, err := server.GetTenantID(ctx)
		if err != nil {
			response.Error(ctx, err)
			ctx.Abort()
			return
		}

		// 检测用户是否为该租户的创建者
		exist, err := orm.Tenants(
			orm.TenantWhere.ID.EQ(tenantID),
			orm.TenantWhere.CreatorID.EQ(userID),
		).ExistsG()
		if err != nil || !exist {
			response.Error(ctx, codes.ErrTenantNotCreator)
		}

		// 获取请求路径和方法
		obj := ctx.Request.URL.Path
		act := strings.ToLower(ctx.Request.Method)

		ok, err := enforcer.Enforce(tenantAdmin, obj, act)
		if err != nil {
			response.Error(ctx, codes.ErrPermissionDenied)
			ctx.Abort()
			return
		}
		if !ok {
			response.Error(ctx, codes.ErrPermissionDenied)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
