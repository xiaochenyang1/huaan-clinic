package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"huaan-medical/internal/middleware"
	"huaan-medical/internal/service"
	"huaan-medical/pkg/errorcode"
	"huaan-medical/pkg/response"
)

// PatientHandler 就诊人处理器
type PatientHandler struct {
	service *service.PatientService
}

// NewPatientHandler 创建就诊人处理器实例
func NewPatientHandler() *PatientHandler {
	return &PatientHandler{
		service: service.NewPatientService(),
	}
}

// List 就诊人列表
// @Summary 获取就诊人列表
// @Description 获取当前用户的就诊人列表
// @Tags 就诊人管理
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} response.Response{data=[]model.PatientVO}
// @Router /api/user/patients [get]
func (h *PatientHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Fail(c, errorcode.ErrUnauthorized)
		return
	}

	list, err := h.service.List(userID)
	if err != nil {
		response.FailWithError(c, err)
		return
	}

	response.Success(c, list)
}

// GetByID 获取就诊人详情
// @Summary 获取就诊人详情
// @Description 根据ID获取就诊人详情
// @Tags 就诊人管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "就诊人ID"
// @Success 200 {object} response.Response{data=model.PatientVO}
// @Router /api/user/patients/{id} [get]
func (h *PatientHandler) GetByID(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Fail(c, errorcode.ErrUnauthorized)
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, errorcode.ErrInvalidIDFormat)
		return
	}

	patient, err := h.service.GetByID(userID, id)
	if err != nil {
		response.FailWithError(c, err)
		return
	}

	response.Success(c, patient)
}

// Create 创建就诊人
// @Summary 创建就诊人
// @Description 添加新的就诊人
// @Tags 就诊人管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body service.CreatePatientRequest true "就诊人信息"
// @Success 200 {object} response.Response{data=model.PatientVO}
// @Router /api/user/patients [post]
func (h *PatientHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Fail(c, errorcode.ErrUnauthorized)
		return
	}

	var req service.CreatePatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errorcode.ErrBindJSON)
		return
	}

	patient, err := h.service.Create(userID, &req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}

	response.Success(c, patient)
}

// Update 更新就诊人
// @Summary 更新就诊人
// @Description 更新就诊人信息
// @Tags 就诊人管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "就诊人ID"
// @Param request body service.UpdatePatientRequest true "就诊人信息"
// @Success 200 {object} response.Response{data=model.PatientVO}
// @Router /api/user/patients/{id} [put]
func (h *PatientHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Fail(c, errorcode.ErrUnauthorized)
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, errorcode.ErrInvalidIDFormat)
		return
	}

	var req service.UpdatePatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errorcode.ErrBindJSON)
		return
	}

	patient, err := h.service.Update(userID, id, &req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}

	response.Success(c, patient)
}

// Delete 删除就诊人
// @Summary 删除就诊人
// @Description 删除就诊人（软删除）
// @Tags 就诊人管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "就诊人ID"
// @Success 200 {object} response.Response
// @Router /api/user/patients/{id} [delete]
func (h *PatientHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Fail(c, errorcode.ErrUnauthorized)
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, errorcode.ErrInvalidIDFormat)
		return
	}

	if err := h.service.Delete(userID, id); err != nil {
		response.FailWithError(c, err)
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

// ListAdmin 查询患者列表（管理后台）
// @Summary 查询患者列表（管理后台）
// @Description 分页查询所有患者，支持关键词搜索
// @Tags 患者管理（后台）
// @Accept json
// @Produce json
// @Security BearerAdmin
// @Param page query int true "页码" minimum(1)
// @Param page_size query int true "每页数量" minimum(1) maximum(100)
// @Param keyword query string false "搜索关键词（姓名、手机号）"
// @Success 200 {object} response.Response{data=response.PageData{list=[]model.PatientVO}}
// @Router /api/admin/patients [get]
func (h *PatientHandler) ListAdmin(c *gin.Context) {
	var req service.ListAdminPatientsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, errorcode.ErrBindJSON)
		return
	}

	patients, total, err := h.service.ListAdmin(&req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}

	response.SuccessWithPage(c, patients, total, req.Page, req.PageSize)
}

// GetByIDAdmin 获取患者详情（管理后台）
// @Summary 获取患者详情（管理后台）
// @Description 获取指定患者的详细信息
// @Tags 患者管理（后台）
// @Accept json
// @Produce json
// @Security BearerAdmin
// @Param id path int true "患者ID"
// @Success 200 {object} response.Response{data=model.PatientVO}
// @Router /api/admin/patients/{id} [get]
func (h *PatientHandler) GetByIDAdmin(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, errorcode.ErrInvalidIDFormat)
		return
	}

	patient, err := h.service.GetByIDAdmin(id)
	if err != nil {
		response.FailWithError(c, err)
		return
	}

	response.Success(c, patient)
}
