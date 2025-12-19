package handler

import (
	"github.com/gin-gonic/gin"

	"huaan-medical/internal/middleware"
	"huaan-medical/internal/service"
	"huaan-medical/pkg/errorcode"
	"huaan-medical/pkg/response"
)

// AdminHandler 管理员处理器
type AdminHandler struct {
	service *service.AdminService
}

// NewAdminHandler 创建管理员处理器实例
func NewAdminHandler() *AdminHandler {
	return &AdminHandler{
		service: service.NewAdminService(),
	}
}

// Login 管理员登录
// @Summary 管理员登录
// @Description 管理员账号密码登录
// @Tags 管理员
// @Accept json
// @Produce json
// @Param request body service.LoginRequest true "登录信息"
// @Success 200 {object} response.Response{data=service.LoginResponse}
// @Router /api/admin/login [post]
func (h *AdminHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errorcode.ErrBindJSON)
		return
	}

	result, err := h.service.Login(&req, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		response.FailWithError(c, err)
		return
	}

	response.Success(c, result)
}

// GetInfo 获取管理员信息
// @Summary 获取当前管理员信息
// @Description 获取当前登录管理员的详细信息
// @Tags 管理员
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} response.Response{data=model.AdminVO}
// @Router /api/admin/info [get]
func (h *AdminHandler) GetInfo(c *gin.Context) {
	adminID := middleware.GetAdminID(c)
	if adminID == 0 {
		response.Fail(c, errorcode.ErrUnauthorized)
		return
	}

	admin, err := h.service.GetAdminInfo(adminID)
	if err != nil {
		response.FailWithError(c, err)
		return
	}

	response.Success(c, admin)
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 修改当前管理员密码
// @Tags 管理员
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body service.ChangePasswordRequest true "密码信息"
// @Success 200 {object} response.Response
// @Router /api/admin/password [put]
func (h *AdminHandler) ChangePassword(c *gin.Context) {
	adminID := middleware.GetAdminID(c)
	if adminID == 0 {
		response.Fail(c, errorcode.ErrUnauthorized)
		return
	}

	var req service.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errorcode.ErrBindJSON)
		return
	}

	if err := h.service.ChangePassword(adminID, &req); err != nil {
		response.FailWithError(c, err)
		return
	}

	response.SuccessWithMessage(c, "密码修改成功", nil)
}
