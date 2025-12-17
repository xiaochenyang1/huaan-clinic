package handler

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"huaan-medical/internal/middleware"
	"huaan-medical/internal/service"
	"huaan-medical/pkg/errorcode"
	"huaan-medical/pkg/response"
)

// AppointmentHandler 预约处理器
type AppointmentHandler struct {
	service *service.AppointmentService
}

// NewAppointmentHandler 创建预约处理器实例
func NewAppointmentHandler() *AppointmentHandler {
	return &AppointmentHandler{
		service: service.NewAppointmentService(),
	}
}

// Create 创建预约
// @Summary 创建预约
// @Description 用户创建新预约
// @Tags 预约
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body service.CreateAppointmentRequest true "预约信息"
// @Success 200 {object} response.Response{data=model.AppointmentVO}
// @Router /api/appointments [post]
func (h *AppointmentHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Fail(c, errorcode.ErrUnauthorized)
		return
	}

	var req service.CreateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errorcode.ErrBindJSON)
		return
	}

	appointment, err := h.service.Create(userID, &req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}

	response.Success(c, appointment)
}

// List 查询用户预约列表
// @Summary 查询用户预约列表
// @Description 查询当前用户的所有预约
// @Tags 预约
// @Accept json
// @Produce json
// @Security Bearer
// @Param status query string false "预约状态（pending/checked_in/completed/cancelled/missed）"
// @Success 200 {object} response.Response{data=[]model.AppointmentListVO}
// @Router /api/appointments [get]
func (h *AppointmentHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Fail(c, errorcode.ErrUnauthorized)
		return
	}

	var req service.ListAppointmentRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, errorcode.ErrBindJSON)
		return
	}

	appointments, err := h.service.ListByUser(userID, &req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}

	response.Success(c, appointments)
}

// GetByID 获取预约详情
// @Summary 获取预约详情
// @Description 获取指定预约的详细信息
// @Tags 预约
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "预约ID"
// @Success 200 {object} response.Response{data=model.AppointmentVO}
// @Router /api/appointments/{id} [get]
func (h *AppointmentHandler) GetByID(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Fail(c, errorcode.ErrUnauthorized)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		response.FailWithMessage(c, errorcode.ErrInvalidParams, "预约ID格式错误")
		return
	}

	appointment, err := h.service.GetByID(userID, id)
	if err != nil {
		response.FailWithError(c, err)
		return
	}

	response.Success(c, appointment)
}

// Cancel 取消预约
// @Summary 取消预约
// @Description 用户取消预约（就诊当天不可取消）
// @Tags 预约
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "预约ID"
// @Param request body service.CancelAppointmentRequest true "取消原因"
// @Success 200 {object} response.Response
// @Router /api/appointments/{id}/cancel [put]
func (h *AppointmentHandler) Cancel(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Fail(c, errorcode.ErrUnauthorized)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		response.FailWithMessage(c, errorcode.ErrInvalidParams, "预约ID格式错误")
		return
	}

	var req service.CancelAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errorcode.ErrBindJSON)
		return
	}

	if err := h.service.Cancel(userID, id, &req); err != nil {
		response.FailWithError(c, err)
		return
	}

	response.Success(c, nil)
}

// Checkin 预约签到
// @Summary 预约签到
// @Description 用户在预约时间范围内签到
// @Tags 预约
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "预约ID"
// @Success 200 {object} response.Response
// @Router /api/appointments/{id}/checkin [post]
func (h *AppointmentHandler) Checkin(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Fail(c, errorcode.ErrUnauthorized)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		response.FailWithMessage(c, errorcode.ErrInvalidParams, "预约ID格式错误")
		return
	}

	if err := h.service.Checkin(userID, id); err != nil {
		response.FailWithError(c, err)
		return
	}

	response.Success(c, nil)
}

// ListAdmin 查询预约列表（管理后台）
// @Summary 查询预约列表（管理后台）
// @Description 分页查询所有预约，支持日期范围、状态、关键词筛选
// @Tags 预约管理
// @Accept json
// @Produce json
// @Security BearerAdmin
// @Param page query int true "页码" minimum(1)
// @Param page_size query int true "每页数量" minimum(1) maximum(100)
// @Param start_date query string false "开始日期（YYYY-MM-DD）"
// @Param end_date query string false "结束日期（YYYY-MM-DD）"
// @Param status query string false "预约状态"
// @Param keyword query string false "搜索关键词（预约编号、患者姓名、医生姓名）"
// @Success 200 {object} response.Response{data=response.PageData{list=[]model.AppointmentVO}}
// @Router /api/admin/appointments [get]
func (h *AppointmentHandler) ListAdmin(c *gin.Context) {
	var req service.ListAdminAppointmentRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, errorcode.ErrBindJSON)
		return
	}

	appointments, total, err := h.service.List(&req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}

	response.SuccessWithPage(c, appointments, total, req.Page, req.PageSize)
}

// UpdateStatus 更新预约状态（管理后台）
// @Summary 更新预约状态（管理后台）
// @Description 管理员更新预约状态
// @Tags 预约管理
// @Accept json
// @Produce json
// @Security BearerAdmin
// @Param id path int true "预约ID"
// @Param request body service.UpdateAppointmentStatusRequest true "状态信息"
// @Success 200 {object} response.Response
// @Router /api/admin/appointments/{id} [put]
func (h *AppointmentHandler) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		response.FailWithMessage(c, errorcode.ErrInvalidParams, "预约ID格式错误")
		return
	}

	var req service.UpdateAppointmentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errorcode.ErrBindJSON)
		return
	}

	if err := h.service.UpdateStatus(id, &req); err != nil {
		response.FailWithError(c, err)
		return
	}

	response.SuccessWithMessage(c, "状态更新成功", nil)
}

// ExportAppointments 导出预约数据（管理后台）
// @Summary 导出预约数据（管理后台）
// @Description 导出预约数据为CSV文件
// @Tags 预约管理
// @Accept json
// @Produce text/csv
// @Security BearerAdmin
// @Param start_date query string false "开始日期（YYYY-MM-DD）"
// @Param end_date query string false "结束日期（YYYY-MM-DD）"
// @Param status query string false "预约状态"
// @Success 200 {file} file "CSV文件"
// @Router /api/admin/appointments/export [get]
func (h *AppointmentHandler) ExportAppointments(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	status := c.Query("status")

	// 构造查询请求（不分页，获取所有数据）
	req := service.ListAdminAppointmentRequest{
		Page:      1,
		PageSize:  10000, // 设置一个较大的值
		StartDate: startDate,
		EndDate:   endDate,
		Status:    func() *string { if status != "" { return &status }; return nil }(),
	}

	appointments, _, err := h.service.List(&req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}

	// 设置响应头
	c.Writer.Header().Set("Content-Type", "text/csv; charset=utf-8")
	c.Writer.Header().Set("Content-Disposition", "attachment; filename=appointments.csv")

	// 写入BOM头（让Excel正确识别UTF-8）
	c.Writer.Write([]byte{0xEF, 0xBB, 0xBF})

	// 写入CSV标题
	c.Writer.Write([]byte("预约编号,患者姓名,医生姓名,科室名称,预约日期,时段,号序,状态,创建时间\n"))

	// 写入数据
	for _, apt := range appointments {
		line := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%d,%s,%s\n",
			apt.AppointmentNo,
			apt.PatientName,
			apt.DoctorName,
			apt.DepartmentName,
			apt.AppointmentDate,
			apt.PeriodName,
			apt.SlotNumber,
			apt.StatusName,
			apt.CreatedAt,
		)
		c.Writer.Write([]byte(line))
	}
}
