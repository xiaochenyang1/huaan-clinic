package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"huaan-medical/internal/service"
	"huaan-medical/pkg/errorcode"
	"huaan-medical/pkg/response"
)

// RoleHandler 角色处理器
type RoleHandler struct {
	service *service.RoleService
}

func NewRoleHandler() *RoleHandler {
	return &RoleHandler{service: service.NewRoleService()}
}

// ListRoles 角色列表
// @Summary 角色列表
// @Description 分页查询角色列表
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security BearerAdmin
// @Param page query int true "页码" minimum(1)
// @Param page_size query int true "每页数量" minimum(1) maximum(100)
// @Param keyword query string false "关键词(编码/名称)"
// @Param status query int false "状态 0禁用 1启用"
// @Success 200 {object} response.Response{data=response.PageData{list=[]model.RoleVO}}
// @Router /api/admin/roles [get]
func (h *RoleHandler) ListRoles(c *gin.Context) {
	var req service.ListRolesRequest
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

// CreateRole 创建角色
// @Summary 创建角色
// @Description 创建角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security BearerAdmin
// @Param request body service.CreateRoleRequest true "角色信息"
// @Success 200 {object} response.Response{data=model.RoleVO}
// @Router /api/admin/roles [post]
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req service.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errorcode.ErrBindJSON)
		return
	}

	role, err := h.service.Create(&req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}

	response.Success(c, role)
}

// UpdateRole 更新角色
// @Summary 更新角色
// @Description 更新角色信息
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security BearerAdmin
// @Param id path int true "角色ID"
// @Param request body service.UpdateRoleRequest true "角色信息"
// @Success 200 {object} response.Response{data=model.RoleVO}
// @Router /api/admin/roles/{id} [put]
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, errorcode.ErrInvalidIDFormat)
		return
	}

	var req service.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errorcode.ErrBindJSON)
		return
	}

	role, err := h.service.Update(id, &req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}

	response.Success(c, role)
}

// UpdateRolePermissions 更新角色权限
// @Summary 更新角色权限
// @Description 替换角色权限（permission_ids 为空则清空）
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security BearerAdmin
// @Param id path int true "角色ID"
// @Param request body service.UpdateRolePermissionsRequest true "权限信息"
// @Success 200 {object} response.Response
// @Router /api/admin/roles/{id}/permissions [put]
func (h *RoleHandler) UpdateRolePermissions(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, errorcode.ErrInvalidIDFormat)
		return
	}

	var req service.UpdateRolePermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errorcode.ErrBindJSON)
		return
	}

	if err := h.service.UpdatePermissions(id, &req); err != nil {
		response.FailWithError(c, err)
		return
	}

	response.SuccessWithMessage(c, "权限更新成功", nil)
}
