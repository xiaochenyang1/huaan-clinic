package handler

import (
	"github.com/gin-gonic/gin"

	"huaan-medical/internal/middleware"
	"huaan-medical/internal/rbac"
	"huaan-medical/internal/service"
	"huaan-medical/pkg/errorcode"
	"huaan-medical/pkg/response"
)

// PermissionHandler 权限处理器
type PermissionHandler struct {
	service *service.PermissionService
}

func NewPermissionHandler() *PermissionHandler {
	return &PermissionHandler{service: service.NewPermissionService()}
}

// ListAllPermissions 权限清单
// @Summary 权限清单
// @Description 查询全部权限（用于角色分配）
// @Tags 权限管理
// @Accept json
// @Produce json
// @Security BearerAdmin
// @Success 200 {object} response.Response{data=[]model.PermissionVO}
// @Router /api/admin/permissions [get]
func (h *PermissionHandler) ListAllPermissions(c *gin.Context) {
	list, err := h.service.ListAll()
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.Success(c, list)
}

type MyPermissionsResponse struct {
	IsSuperAdmin bool     `json:"is_super_admin"`
	Permissions  []string `json:"permissions"`
}

// MyPermissions 当前管理员权限
// @Summary 当前管理员权限
// @Description 获取当前登录管理员的权限码列表（用于前端菜单/按钮权限）
// @Tags 权限管理
// @Accept json
// @Produce json
// @Security BearerAdmin
// @Success 200 {object} response.Response{data=handler.MyPermissionsResponse}
// @Router /api/admin/permissions/me [get]
func (h *PermissionHandler) MyPermissions(c *gin.Context) {
	adminID := middleware.GetAdminID(c)
	if adminID == 0 {
		response.Fail(c, errorcode.ErrUnauthorized)
		return
	}

	codes, isSuper, err := h.service.GetAdminPermissionCodes(adminID)
	if err != nil {
		response.FailWithError(c, err)
		return
	}

	if isSuper {
		set := map[string]struct{}{}
		for _, code := range codes {
			set[code] = struct{}{}
		}
		for _, p := range rbac.DefaultPermissions {
			set[p.Code] = struct{}{}
		}
		codes = make([]string, 0, len(set))
		for code := range set {
			codes = append(codes, code)
		}
	}

	response.Success(c, MyPermissionsResponse{
		IsSuperAdmin: isSuper,
		Permissions:  codes,
	})
}
