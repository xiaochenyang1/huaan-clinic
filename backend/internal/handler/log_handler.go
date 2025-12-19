package handler

import (
	"github.com/gin-gonic/gin"

	"huaan-medical/internal/service"
	"huaan-medical/pkg/errorcode"
	"huaan-medical/pkg/response"
)

// LogHandler 日志处理器
type LogHandler struct {
	service *service.LogService
}

func NewLogHandler() *LogHandler {
	return &LogHandler{service: service.NewLogService()}
}

// ListOperationLogs 操作日志列表（管理后台）
// @Summary 查询操作日志（管理后台）
// @Description 分页查询操作日志，支持关键词/模块/日期筛选
// @Tags 操作日志
// @Accept json
// @Produce json
// @Security BearerAdmin
// @Param page query int true "页码" minimum(1)
// @Param page_size query int true "每页数量" minimum(1) maximum(100)
// @Param admin_id query int false "管理员ID"
// @Param module query string false "模块"
// @Param keyword query string false "关键词"
// @Param start_date query string false "开始日期（YYYY-MM-DD）"
// @Param end_date query string false "结束日期（YYYY-MM-DD）"
// @Success 200 {object} response.Response{data=response.PageData{list=[]model.OperationLogVO}}
// @Router /api/admin/logs/operation [get]
func (h *LogHandler) ListOperationLogs(c *gin.Context) {
	var req service.ListOperationLogsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, errorcode.ErrBindJSON)
		return
	}

	list, total, err := h.service.ListOperationLogs(&req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}

	response.SuccessWithPage(c, list, total, req.Page, req.PageSize)
}

// ListLoginLogs 登录日志列表（管理后台）
// @Summary 查询登录日志（管理后台）
// @Description 分页查询登录日志，支持关键词/状态/日期筛选
// @Tags 登录日志
// @Accept json
// @Produce json
// @Security BearerAdmin
// @Param page query int true "页码" minimum(1)
// @Param page_size query int true "每页数量" minimum(1) maximum(100)
// @Param user_type query string false "用户类型 user/admin"
// @Param status query int false "状态 0失败 1成功"
// @Param keyword query string false "关键词"
// @Param start_date query string false "开始日期（YYYY-MM-DD）"
// @Param end_date query string false "结束日期（YYYY-MM-DD）"
// @Success 200 {object} response.Response{data=response.PageData{list=[]model.LoginLogVO}}
// @Router /api/admin/logs/login [get]
func (h *LogHandler) ListLoginLogs(c *gin.Context) {
	var req service.ListLoginLogsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, errorcode.ErrBindJSON)
		return
	}

	list, total, err := h.service.ListLoginLogs(&req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}

	response.SuccessWithPage(c, list, total, req.Page, req.PageSize)
}

