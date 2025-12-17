package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"huaan-medical/internal/service"
	"huaan-medical/pkg/errorcode"
	"huaan-medical/pkg/response"
)

// DoctorHandler 医生处理器
type DoctorHandler struct {
	service *service.DoctorService
}

// NewDoctorHandler 创建医生处理器实例
func NewDoctorHandler() *DoctorHandler {
	return &DoctorHandler{
		service: service.NewDoctorService(),
	}
}

// List 医生列表（管理后台）
// @Summary 医生列表
// @Description 分页查询医生列表（管理后台）
// @Tags 医生管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int true "页码"
// @Param page_size query int true "每页数量"
// @Param department_id query int false "科室ID筛选"
// @Param status query int false "状态筛选"
// @Param keyword query string false "关键词搜索"
// @Success 200 {object} response.Response{data=response.PageData}
// @Router /api/admin/doctors [get]
func (h *DoctorHandler) List(c *gin.Context) {
	var req service.ListDoctorRequest
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

// ListPublic 医生列表（公开接口）
// @Summary 获取医生列表
// @Description 获取启用的医生列表（公开接口）
// @Tags 医生
// @Accept json
// @Produce json
// @Param department_id query int false "科室ID筛选"
// @Param keyword query string false "关键词搜索"
// @Success 200 {object} response.Response{data=[]model.DoctorListVO}
// @Router /api/doctors [get]
func (h *DoctorHandler) ListPublic(c *gin.Context) {
	var req service.ListPublicDoctorRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, errorcode.ErrInvalidParams)
		return
	}

	list, err := h.service.ListPublic(&req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}

	response.Success(c, list)
}

// GetByID 获取医生详情
// @Summary 获取医生详情
// @Description 根据ID获取医生详情
// @Tags 医生管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "医生ID"
// @Success 200 {object} response.Response{data=model.DoctorVO}
// @Router /api/admin/doctors/{id} [get]
func (h *DoctorHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, errorcode.ErrInvalidIDFormat)
		return
	}

	doctor, err := h.service.GetByID(id)
	if err != nil {
		response.FailWithError(c, err)
		return
	}

	response.Success(c, doctor)
}

// GetByIDPublic 获取医生详情（公开接口）
// @Summary 获取医生详情
// @Description 根据ID获取医生详情（公开接口）
// @Tags 医生
// @Accept json
// @Produce json
// @Param id path int true "医生ID"
// @Success 200 {object} response.Response{data=model.DoctorVO}
// @Router /api/doctors/{id} [get]
func (h *DoctorHandler) GetByIDPublic(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, errorcode.ErrInvalidIDFormat)
		return
	}

	doctor, err := h.service.GetByID(id)
	if err != nil {
		response.FailWithError(c, err)
		return
	}

	response.Success(c, doctor)
}

// Create 创建医生
// @Summary 创建医生
// @Description 创建新医生
// @Tags 医生管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body service.CreateDoctorRequest true "医生信息"
// @Success 200 {object} response.Response{data=model.DoctorVO}
// @Router /api/admin/doctors [post]
func (h *DoctorHandler) Create(c *gin.Context) {
	var req service.CreateDoctorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errorcode.ErrBindJSON)
		return
	}

	doctor, err := h.service.Create(&req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}

	response.Success(c, doctor)
}

// Update 更新医生
// @Summary 更新医生
// @Description 更新医生信息
// @Tags 医生管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "医生ID"
// @Param request body service.UpdateDoctorRequest true "医生信息"
// @Success 200 {object} response.Response{data=model.DoctorVO}
// @Router /api/admin/doctors/{id} [put]
func (h *DoctorHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, errorcode.ErrInvalidIDFormat)
		return
	}

	var req service.UpdateDoctorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errorcode.ErrBindJSON)
		return
	}

	doctor, err := h.service.Update(id, &req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}

	response.Success(c, doctor)
}

// Delete 删除医生
// @Summary 删除医生
// @Description 删除医生（软删除）
// @Tags 医生管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "医生ID"
// @Success 200 {object} response.Response
// @Router /api/admin/doctors/{id} [delete]
func (h *DoctorHandler) Delete(c *gin.Context) {
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
