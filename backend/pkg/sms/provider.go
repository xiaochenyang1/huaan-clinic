package sms

import "context"

// Provider 短信服务商抽象
// 说明：仅负责把验证码发送给目标手机号；验证码生成/频控/存储由 Service 负责。
type Provider interface {
	SendCode(ctx context.Context, phone string, code string) error
}
