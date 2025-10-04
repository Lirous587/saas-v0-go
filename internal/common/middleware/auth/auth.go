package auth

import (
	casbinadapter "saas/internal/common/casbin"
	"saas/internal/common/reskit/codes"
	"saas/internal/common/reskit/response"
	"saas/internal/common/server"
	roleadapter "saas/internal/role/adapters"
	roleDomain "saas/internal/role/domain"
	roleService "saas/internal/role/service"
	useradapter "saas/internal/user/adapters"
	userdomain "saas/internal/user/domain"
	userService "saas/internal/user/service"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
)

var tokenServer userdomain.TokenService

var roleServer roleDomain.RoleService

var enforcer *casbin.Enforcer

func Init() {
	// 初始化token服务
	tokenCache := useradapter.NewTokenRedisCache()
	userRepo := useradapter.NewUserPSQLRepository()
	tokenServer = userService.NewTokenService(tokenCache, userRepo)

	// 初始化roleService
	roleCache := roleadapter.NewRoleRedisCache()
	roleRepo := roleadapter.NewRolePSQLRepository()
	roleServer = roleService.NewRoleService(roleRepo, roleCache)

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

		// 查询用户在该租户下的角色
		role, err := roleServer.GetUserRoleInTenant(userID, tenantID)
		if err != nil {
			response.Error(ctx, codes.ErrRoleNotFound)
			ctx.Abort()
			return
		}

		// 获取请求路径和方法
		obj := ctx.Request.URL.Path
		act := strings.ToLower(ctx.Request.Method)

		// 将 tenantID 转换为字符串（确保与策略类型匹配）
		// tenantIDStr := strconv.FormatInt(tenantID, 10)

		ok, err := enforcer.Enforce(role.Name, obj, act)
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
