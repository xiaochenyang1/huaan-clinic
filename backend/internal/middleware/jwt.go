package middleware

import (
	"context"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"

	"huaan-medical/pkg/errorcode"
	"huaan-medical/pkg/jwt"
	"huaan-medical/pkg/redis"
	"huaan-medical/pkg/response"
)

const (
	// ContextKeyUserID 用户ID上下文键
	ContextKeyUserID = "user_id"
	// ContextKeyOpenID OpenID上下文键
	ContextKeyOpenID = "open_id"
	// ContextKeyAdminID 管理员ID上下文键
	ContextKeyAdminID = "admin_id"
	// ContextKeyAdminRole 管理员角色上下文键
	ContextKeyAdminRole = "admin_role"
)

// JWTAuth 用户JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			response.Fail(c, errorcode.ErrUnauthorized)
			c.Abort()
			return
		}

		// 检查Token是否在黑名单中（需要Redis）
		if redis.IsEnabled() {
			blacklistKey := fmt.Sprintf(redis.KeyTokenBlacklist, token)
			exists, _ := redis.Exists(context.Background(), blacklistKey)
			if exists {
				response.Fail(c, errorcode.ErrTokenInvalid)
				c.Abort()
				return
			}
		}

		// 解析Token
		claims, err := jwt.ParseToken(token)
		if err != nil {
			if strings.Contains(err.Error(), "过期") {
				response.Fail(c, errorcode.ErrTokenExpired)
			} else {
				response.Fail(c, errorcode.ErrTokenInvalid)
			}
			c.Abort()
			return
		}

		// 验证Token类型
		if claims.TokenType != jwt.AccessToken {
			response.Fail(c, errorcode.ErrTokenInvalid)
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyOpenID, claims.OpenID)
		c.Next()
	}
}

// JWTAdminAuth 管理员JWT认证中间件
func JWTAdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			response.Fail(c, errorcode.ErrUnauthorized)
			c.Abort()
			return
		}

		// 检查Token是否在黑名单中（需要Redis）
		if redis.IsEnabled() {
			blacklistKey := fmt.Sprintf(redis.KeyTokenBlacklist, token)
			exists, _ := redis.Exists(context.Background(), blacklistKey)
			if exists {
				response.Fail(c, errorcode.ErrTokenInvalid)
				c.Abort()
				return
			}
		}

		// 解析Token
		claims, err := jwt.ParseAdminToken(token)
		if err != nil {
			if strings.Contains(err.Error(), "过期") {
				response.Fail(c, errorcode.ErrTokenExpired)
			} else {
				response.Fail(c, errorcode.ErrTokenInvalid)
			}
			c.Abort()
			return
		}

		// 将管理员信息存入上下文
		c.Set(ContextKeyAdminID, claims.AdminID)
		c.Set(ContextKeyAdminRole, claims.Role)
		c.Next()
	}
}

// extractToken 从请求头中提取Token
func extractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	// Bearer Token格式
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

// GetUserID 从上下文获取��户ID
func GetUserID(c *gin.Context) int64 {
	if userID, exists := c.Get(ContextKeyUserID); exists {
		return userID.(int64)
	}
	return 0
}

// GetOpenID 从上下文获取OpenID
func GetOpenID(c *gin.Context) string {
	if openID, exists := c.Get(ContextKeyOpenID); exists {
		return openID.(string)
	}
	return ""
}

// GetAdminID 从上下文获取管理员ID
func GetAdminID(c *gin.Context) int64 {
	if adminID, exists := c.Get(ContextKeyAdminID); exists {
		return adminID.(int64)
	}
	return 0
}

// GetAdminRole 从上下文获取管理员角色
func GetAdminRole(c *gin.Context) string {
	if role, exists := c.Get(ContextKeyAdminRole); exists {
		return role.(string)
	}
	return ""
}
