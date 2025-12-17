package handler

import (
	"github.com/gin-gonic/gin"

	"huaan-medical/internal/service"
	"huaan-medical/pkg/response"
)

// StatisticsHandler 统计处理器
type StatisticsHandler struct {
	service *service.StatisticsService
}

// NewStatisticsHandler 创建统计处理器实例
func NewStatisticsHandler() *StatisticsHandler {
	return &StatisticsHandler{
		service: service.NewStatisticsService(),
	}
}

// GetDashboard 获取仪表盘数据
// @Summary 获取仪表盘数据
// @Description 获取管理后台仪表盘统计数据
// @Tags 数据统计
// @Accept json
// @Produce json
// @Security BearerAdmin
// @Success 200 {object} response.Response{data=service.DashboardData}
// @Router /api/admin/dashboard [get]
func (h *StatisticsHandler) GetDashboard(c *gin.Context) {
	data, err := h.service.GetDashboard()
	if err != nil {
		response.FailWithError(c, err)
		return
	}

	response.Success(c, data)
}

// GetStatistics 获取统计数据
// @Summary 获取统计数据
// @Description 获取各类统计数据，支持日期范围筛选
// @Tags 数据统计
// @Accept json
// @Produce json
// @Security BearerAdmin
// @Param start_date query string false "开始日期（YYYY-MM-DD）"
// @Param end_date query string false "结束日期（YYYY-MM-DD）"
// @Success 200 {object} response.Response{data=service.StatisticsData}
// @Router /api/admin/statistics [get]
func (h *StatisticsHandler) GetStatistics(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	data, err := h.service.GetStatistics(startDate, endDate)
	if err != nil {
		response.FailWithError(c, err)
		return
	}

	response.Success(c, data)
}
