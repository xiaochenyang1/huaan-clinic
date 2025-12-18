package handler

import (
	"github.com/gin-gonic/gin"

	"huaan-medical/pkg/config"
	"huaan-medical/pkg/errorcode"
	"huaan-medical/pkg/response"
	"huaan-medical/pkg/sms"
	"huaan-medical/pkg/utils"
)

// SMSHandler 短信处理器
type SMSHandler struct{}

// NewSMSHandler 创建短信处理器实例
func NewSMSHandler() *SMSHandler {
	return &SMSHandler{}
}

// SendCodeRequest 发送验证码请求
type SendCodeRequest struct {
	Phone string `json:"phone" binding:"required,len=11"`
}

// SendCodeResponse 发送验证码响应
type SendCodeResponse struct {
	Code string `json:"code,omitempty"` // 仅测试环境返回
}

// SendCode 发送验证码
// @Summary 发送验证码
// @Description 发送短信验证码（测试环境返回验证码）
// @Tags 短信
// @Accept json
// @Produce json
// @Param request body SendCodeRequest true "手机号"
// @Success 200 {object} response.Response{data=SendCodeResponse}
// @Router /api/sms/send [post]
func (h *SMSHandler) SendCode(c *gin.Context) {
	var req SendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errorcode.ErrBindJSON)
		return
	}

	cfg, err := config.Load("config.yaml")
	if err != nil {
		response.FailWithMessage(c, errorcode.ErrInternalServer, "配置加载失败")
		return
	}
	if !cfg.SMS.Enabled {
		response.FailWithMessage(c, errorcode.ErrInvalidParams, "短信服务未启用")
		return
	}

	// 验证手机号格式
	if !utils.ValidatePhone(req.Phone) {
		response.FailWithMessage(c, errorcode.ErrInvalidParams, "手机号格式错误")
		return
	}

	// 发送验证码
	smsService := sms.GetService()
	code, err := smsService.SendCode(req.Phone)
	if err != nil {
		response.FailWithMessage(c, errorcode.ErrSMSCodeSendTooFrequent, err.Error())
		return
	}

	// 测试环境返回验证码，生产环境不返回
	resp := &SendCodeResponse{}
	if cfg.SMS.AllowTestCode && gin.Mode() != gin.ReleaseMode {
		resp.Code = code // 仅测试/开发环境返回
	}

	response.SuccessWithMessage(c, "验证码发送成功", resp)
}
