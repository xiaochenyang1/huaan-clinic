package handler

import (
	"github.com/gin-gonic/gin"

	"huaan-medical/internal/middleware"
	"huaan-medical/internal/service"
	"huaan-medical/pkg/errorcode"
	"huaan-medical/pkg/response"
)

// UserHandler 用户处理器
type UserHandler struct {
	service *service.UserService
}

// NewUserHandler 创建用户处理器实例
func NewUserHandler() *UserHandler {
	return &UserHandler{
		service: service.NewUserService(),
	}
}

// WeChatLogin 微信登录
// @Summary 微信登录
// @Description 微信小程序登录
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body service.WeChatLoginRequest true "登录信息"
// @Success 200 {object} response.Response{data=service.WeChatLoginResponse}
// @Router /api/user/login [post]
func (h *UserHandler) WeChatLogin(c *gin.Context) {
	var req service.WeChatLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errorcode.ErrBindJSON)
		return
	}

	result, err := h.service.WeChatLogin(&req, c.ClientIP())
	if err != nil {
		response.FailWithError(c, err)
		return
	}

	response.Success(c, result)
}

// GetInfo 获取用户信息
// @Summary 获取用户信息
// @Description 获取当前登录用户信息
// @Tags 用户
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} response.Response{data=model.UserVO}
// @Router /api/user/info [get]
func (h *UserHandler) GetInfo(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Fail(c, errorcode.ErrUnauthorized)
		return
	}

	user, err := h.service.GetUserInfo(userID)
	if err != nil {
		response.FailWithError(c, err)
		return
	}

	response.Success(c, user)
}

// UpdateInfo 更新用户信息
// @Summary 更新用户信息
// @Description 更新当前用户信息
// @Tags 用户
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body service.UpdateUserInfoRequest true "用户信息"
// @Success 200 {object} response.Response{data=model.UserVO}
// @Router /api/user/info [put]
func (h *UserHandler) UpdateInfo(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Fail(c, errorcode.ErrUnauthorized)
		return
	}

	var req service.UpdateUserInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errorcode.ErrBindJSON)
		return
	}

	user, err := h.service.UpdateUserInfo(userID, &req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}

	response.Success(c, user)
}
