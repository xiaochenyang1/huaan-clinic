package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"huaan-medical/internal/middleware"
	"huaan-medical/internal/service"
	"huaan-medical/pkg/errorcode"
	"huaan-medical/pkg/response"
)

// MedicalRecordHandler 就诊记录处理器
type MedicalRecordHandler struct {
	service *service.MedicalRecordService
}

// NewMedicalRecordHandler 创建就诊记录处理器实例
func NewMedicalRecordHandler() *MedicalRecordHandler {
	return &MedicalRecordHandler{
		service: service.NewMedicalRecordService(),
	}
}

// List 查询就诊记录列表
// @Summary 查询就诊记录列表
// @Description 查询当前用户的所有就诊记录
// @Tags 就诊记录
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} response.Response{data=[]model.MedicalRecordListVO}
// @Router /api/records [get]
func (h *MedicalRecordHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Fail(c, errorcode.ErrUnauthorized)
		return
	}

	records, err := h.service.ListByUser(userID)
	if err != nil {
		response.FailWithError(c, err)
		return
	}

	response.Success(c, records)
}

// GetByID 获取就诊记录详情
// @Summary 获取就诊记录详情
// @Description 获取指定就诊记录的详细信息
// @Tags 就诊记录
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "就诊记录ID"
// @Success 200 {object} response.Response{data=model.MedicalRecordVO}
// @Router /api/records/{id} [get]
func (h *MedicalRecordHandler) GetByID(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Fail(c, errorcode.ErrUnauthorized)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		response.FailWithMessage(c, errorcode.ErrInvalidParams, "就诊记录ID格式错误")
		return
	}

	record, err := h.service.GetByID(userID, id)
	if err != nil {
		response.FailWithError(c, err)
		return
	}

	response.Success(c, record)
}
