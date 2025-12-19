package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"huaan-medical/internal/service"
	"huaan-medical/pkg/errorcode"
	"huaan-medical/pkg/response"
)

// AdminManageHandler 管理员管理处理器
type AdminManageHandler struct {
	service *service.AdminManageService
}

func NewAdminManageHandler() *AdminManageHandler {
	return &AdminManageHandler{service: service.NewAdminManageService()}
}

// ListAdmins 管理员列表
// @Summary 管理员列表
// @Description 分页查询管理员列表
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Security BearerAdmin
// @Param page query int true "页码" minimum(1)
// @Param page_size query int true "每页数量" minimum(1) maximum(100)
// @Param keyword query string false "关键词(用户名/昵称/手机号/邮箱)"
// @Param status query int false "状态 0禁用 1启用"
// @Success 200 {object} response.Response{data=response.PageData{list=[]model.AdminVO}}
// @Router /api/admin/admins [get]
func (h *AdminManageHandler) ListAdmins(c *gin.Context) {
	var req service.ListAdminsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, errorcode.ErrInvalidPageParams)
		return
	}

	list, total, err := h.service.List(&req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.SuccessWithPage(c, list, total, req.Page, req.PageSize)
}

// CreateAdmin 创建管理员
// @Summary 创建管理员
// @Description 创建管理员并可选绑定角色
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Security BearerAdmin
// @Param request body service.CreateAdminRequest true "管理员信息"
// @Success 200 {object} response.Response{data=model.AdminVO}
// @Router /api/admin/admins [post]
func (h *AdminManageHandler) CreateAdmin(c *gin.Context) {
	var req service.CreateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errorcode.ErrBindJSON)
		return
	}

	admin, err := h.service.Create(&req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}

	response.Success(c, admin)
}

// UpdateAdmin 更新管理员
// @Summary 更新管理员
// @Description 更新管理员信息与角色
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Security BearerAdmin
// @Param id path int true "管理员ID"
// @Param request body service.UpdateAdminRequest true "管理员信息"
// @Success 200 {object} response.Response{data=model.AdminVO}
// @Router /api/admin/admins/{id} [put]
func (h *AdminManageHandler) UpdateAdmin(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, errorcode.ErrInvalidIDFormat)
		return
	}

	var req service.UpdateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errorcode.ErrBindJSON)
		return
	}

	admin, err := h.service.Update(id, &req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}

	response.Success(c, admin)
}

// ResetAdminPassword 重置管理员密码
// @Summary 重置管理员密码
// @Description 重置管理员登录密码
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Security BearerAdmin
// @Param id path int true "管理员ID"
// @Param request body service.ResetAdminPasswordRequest true "密码信息"
// @Success 200 {object} response.Response
// @Router /api/admin/admins/{id}/password [put]
func (h *AdminManageHandler) ResetAdminPassword(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, errorcode.ErrInvalidIDFormat)
		return
	}

	var req service.ResetAdminPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errorcode.ErrBindJSON)
		return
	}

	if err := h.service.ResetPassword(id, &req); err != nil {
		response.FailWithError(c, err)
		return
	}

	response.SuccessWithMessage(c, "密码重置成功", nil)
}
