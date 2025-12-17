package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"huaan-medical/internal/service"
	"huaan-medical/pkg/errorcode"
	"huaan-medical/pkg/response"
)

// ScheduleHandler 排班处理器
type ScheduleHandler struct {
	service *service.ScheduleService
}

// NewScheduleHandler 创建排班处理器实例
func NewScheduleHandler() *ScheduleHandler {
	return &ScheduleHandler{
		service: service.NewScheduleService(),
	}
}

// List 排班列表（管理后台）
// @Summary 排班列表
// @Description 分页查询排班列表（管理后台）
// @Tags 排班管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int true "页码"
// @Param page_size query int true "每页数量"
// @Param doctor_id query int false "医生ID筛选"
// @Param department_id query int false "科室ID筛选"
// @Param start_date query string false "开始日期 YYYY-MM-DD"
// @Param end_date query string false "结束日期 YYYY-MM-DD"
// @Param status query int false "状态筛选"
// @Success 200 {object} response.Response{data=response.PageData}
// @Router /api/admin/schedules [get]
func (h *ScheduleHandler) List(c *gin.Context) {
	var req service.ListScheduleRequest
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

// GetByID 获取排班详情
// @Summary 获取排班详情
// @Description 根据ID获取排班详情
// @Tags 排班管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "排班ID"
// @Success 200 {object} response.Response{data=model.ScheduleVO}
// @Router /api/admin/schedules/{id} [get]
func (h *ScheduleHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, errorcode.ErrInvalidIDFormat)
		return
	}

	schedule, err := h.service.GetByID(id)
	if err != nil {
		response.FailWithError(c, err)
		return
	}

	response.Success(c, schedule)
}

// Create 创建排班
// @Summary 创建排班
// @Description 创建单个排班
// @Tags 排班管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body service.CreateScheduleRequest true "排班信息"
// @Success 200 {object} response.Response{data=model.ScheduleVO}
// @Router /api/admin/schedules [post]
func (h *ScheduleHandler) Create(c *gin.Context) {
	var req service.CreateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errorcode.ErrBindJSON)
		return
	}

	schedule, err := h.service.Create(&req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}

	response.Success(c, schedule)
}

// BatchCreate 批量创建排班
// @Summary 批量创建排班
// @Description 批量创建多天/多时段的排班
// @Tags 排班管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body service.BatchCreateScheduleRequest true "批量排班信息"
// @Success 200 {object} response.Response{data=map[string]int}
// @Router /api/admin/schedules/batch [post]
func (h *ScheduleHandler) BatchCreate(c *gin.Context) {
	var req service.BatchCreateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errorcode.ErrBindJSON)
		return
	}

	count, err := h.service.BatchCreate(&req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}

	response.Success(c, gin.H{
		"created_count": count,
		"message":       "批量创建成功",
	})
}

// Update 更新排班
// @Summary 更新排班
// @Description 更新排班信息
// @Tags 排班管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "排班ID"
// @Param request body service.UpdateScheduleRequest true "排班信息"
// @Success 200 {object} response.Response{data=model.ScheduleVO}
// @Router /api/admin/schedules/{id} [put]
func (h *ScheduleHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, errorcode.ErrInvalidIDFormat)
		return
	}

	var req service.UpdateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errorcode.ErrBindJSON)
		return
	}

	schedule, err := h.service.Update(id, &req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}

	response.Success(c, schedule)
}

// Delete 删除排班
// @Summary 删除排班
// @Description 删除排班（软删除）
// @Tags 排班管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "排班ID"
// @Success 200 {object} response.Response
// @Router /api/admin/schedules/{id} [delete]
func (h *ScheduleHandler) Delete(c *gin.Context) {
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

// ListByDoctor 查询医生排班（公开接口）
// @Summary 查询医生排班
// @Description 查询医生在指定日期范围内的排班
// @Tags 排班
// @Accept json
// @Produce json
// @Param doctor_id query int true "医生ID"
// @Param start_date query string true "开始日期 YYYY-MM-DD"
// @Param end_date query string true "结束日期 YYYY-MM-DD"
// @Success 200 {object} response.Response{data=[]model.ScheduleVO}
// @Router /api/schedule [get]
func (h *ScheduleHandler) ListByDoctor(c *gin.Context) {
	doctorIDStr := c.Query("doctor_id")
	if doctorIDStr == "" {
		response.FailWithMessage(c, errorcode.ErrInvalidParams, "医生ID不能为空")
		return
	}

	doctorID, err := strconv.ParseInt(doctorIDStr, 10, 64)
	if err != nil {
		response.FailWithMessage(c, errorcode.ErrInvalidParams, "医生ID格式错误")
		return
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		response.FailWithMessage(c, errorcode.ErrInvalidParams, "日期范围不能为空")
		return
	}

	list, err := h.service.ListByDoctor(doctorID, startDate, endDate)
	if err != nil {
		response.FailWithError(c, err)
		return
	}

	response.Success(c, list)
}

// ListAvailable 查询可预约时段（公开接口）
// @Summary 获取可预约时段
// @Description 查询有剩余号源的排班列表
// @Tags 排班
// @Accept json
// @Produce json
// @Param doctor_id query int false "医生ID筛选"
// @Param department_id query int false "科室ID筛选"
// @Param start_date query string true "开始日期 YYYY-MM-DD"
// @Param end_date query string true "结束日期 YYYY-MM-DD"
// @Success 200 {object} response.Response{data=[]model.ScheduleVO}
// @Router /api/schedule/available [get]
func (h *ScheduleHandler) ListAvailable(c *gin.Context) {
	var req service.ListAvailableScheduleRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, errorcode.ErrInvalidParams)
		return
	}

	list, err := h.service.ListAvailable(&req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}

	response.Success(c, list)
}
