package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"huaan-medical/internal/model"
	"huaan-medical/internal/rbac"
	"huaan-medical/internal/repository"
	"huaan-medical/pkg/database"
	"huaan-medical/pkg/errorcode"
	"huaan-medical/pkg/response"
)

const (
	// ContextKeyAdminPermissions 管理员权限码上下文键
	ContextKeyAdminPermissions = "admin_permissions"
)

// AdminRBAC 管理后台 RBAC 鉴权中间件
func AdminRBAC() gin.HandlerFunc {
	return func(c *gin.Context) {
		fullPath := c.FullPath()
		if fullPath == "" {
			c.Next()
			return
		}

		key := c.Request.Method + " " + fullPath
		required, exists := rbac.AdminRoutePermissions[key]
		if !exists {
			response.FailWithMessage(c, errorcode.ErrPermissionDenied, "接口未配置权限")
			c.Abort()
			return
		}

		// 不需要权限（登录后即可访问）
		if len(required) == 0 {
			c.Next()
			return
		}

		// 如果权限表为空，默认放行（避免初始化阶段直接锁死后台）
		db := database.GetDB()
		if db != nil {
			var permCount int64
			if err := db.Model(&model.Permission{}).Count(&permCount).Error; err == nil && permCount == 0 {
				c.Next()
				return
			}
		}

		adminID := GetAdminID(c)
		if adminID == 0 {
			response.Fail(c, errorcode.ErrUnauthorized)
			c.Abort()
			return
		}

		adminRepo := repository.NewAdminRepository()
		admin, err := adminRepo.GetByIDWithRolesPermissions(adminID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				response.Fail(c, errorcode.ErrUnauthorized)
			} else {
				response.Fail(c, errorcode.ErrDatabase)
			}
			c.Abort()
			return
		}

		permSet := map[string]struct{}{}
		for _, role := range admin.Roles {
			// 超级管理员默认拥有全部权限
			if role.Code == rbac.RoleSuperAdmin {
				c.Next()
				return
			}
			for _, p := range role.Permissions {
				if p.Code == "" {
					continue
				}
				permSet[p.Code] = struct{}{}
			}
		}

		allowed := false
		for _, p := range required {
			if _, ok := permSet[p]; ok {
				allowed = true
				break
			}
		}
		if !allowed {
			response.Fail(c, errorcode.ErrPermissionDenied)
			c.Abort()
			return
		}

		// 保存到上下文（供后续 handler 使用，可选）
		perms := make([]string, 0, len(permSet))
		for code := range permSet {
			perms = append(perms, code)
		}
		c.Set(ContextKeyAdminPermissions, perms)

		c.Next()
	}
}
