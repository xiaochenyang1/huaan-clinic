package sms

import (
	"context"

	"go.uber.org/zap"

	"huaan-medical/pkg/logger"
)

// ConsoleProvider 开发/联调用：把验证码输出到日志（不回传给客户端）。
type ConsoleProvider struct{}

func (p *ConsoleProvider) SendCode(_ context.Context, phone string, code string) error {
	logger.Warn("短信验证码（console provider）", zap.String("phone", maskPhone(phone)), zap.String("code", code))
	return nil
}

func maskPhone(phone string) string {
	if len(phone) != 11 {
		return phone
	}
	return phone[:3] + "****" + phone[7:]
}
