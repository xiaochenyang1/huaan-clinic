package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"huaan-medical/internal/service"
	"huaan-medical/pkg/errorcode"
	"huaan-medical/pkg/response"
)

// DepartmentHandler 科室处理器
type DepartmentHandler struct {
	service *service.DepartmentService
}

// NewDepartmentHandler 创建科室处理器实例
func NewDepartmentHandler() *DepartmentHandler {
	return &DepartmentHandler{
		service: service.NewDepartmentService(),
	}
}

// List 科室列表（管理后台）
// @Summary 科室列表
// @Description 分页查询科室列表（管理后台）
// @Tags 科室管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int true "页码"
// @Param page_size query int true "每页数量"
// @Param status query int false "状态筛选"
// @Success 200 {object} response.Response{data=response.PageData}
// @Router /api/admin/departments [get]
func (h *DepartmentHandler) List(c *gin.Context) {
	var req service.ListDepartmentRequest
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

// ListAll 科室列表（公开接口）
// @Summary 获取所有科室
// @Description 获取所有启用的科室列表（公开接口）
// @Tags 科室
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]model.DepartmentVO}
// @Router /api/departments [get]
func (h *DepartmentHandler) ListAll(c *gin.Context) {
	list, err := h.service.ListAll()
	if err != nil {
		response.FailWithError(c, err)
		return
	}

	response.Success(c, list)
}

// GetByID 获取科室详情
// @Summary 获取科室详情
// @Description 根据ID获取科室详情
// @Tags 科室管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "科室ID"
// @Success 200 {object} response.Response{data=model.DepartmentVO}
// @Router /api/admin/departments/{id} [get]
func (h *DepartmentHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, errorcode.ErrInvalidIDFormat)
		return
	}

	dept, err := h.service.GetByID(id)
	if err != nil {
		response.FailWithError(c, err)
		return
	}

	response.Success(c, dept)
}

// Create 创建科室
// @Summary 创建科室
// @Description 创建新科室
// @Tags 科室管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body service.CreateDepartmentRequest true "科室信息"
// @Success 200 {object} response.Response{data=model.DepartmentVO}
// @Router /api/admin/departments [post]
func (h *DepartmentHandler) Create(c *gin.Context) {
	var req service.CreateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errorcode.ErrBindJSON)
		return
	}

	dept, err := h.service.Create(&req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}

	response.Success(c, dept)
}

// Update 更新科室
// @Summary 更新科室
// @Description 更新科室信息
// @Tags 科室管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "科室ID"
// @Param request body service.UpdateDepartmentRequest true "科室信息"
// @Success 200 {object} response.Response{data=model.DepartmentVO}
// @Router /api/admin/departments/{id} [put]
func (h *DepartmentHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, errorcode.ErrInvalidIDFormat)
		return
	}

	var req service.UpdateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errorcode.ErrBindJSON)
		return
	}

	dept, err := h.service.Update(id, &req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}

	response.Success(c, dept)
}

// Delete 删除科室
// @Summary 删除科室
// @Description 删除科室（软删除）
// @Tags 科室管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "科室ID"
// @Success 200 {object} response.Response
// @Router /api/admin/departments/{id} [delete]
func (h *DepartmentHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, errorcode.ErrInvalidIDFormat)
		return
	}

	if err := h.service.Delete(id); err != nil {
		response.FailWithError(c, err)
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}
