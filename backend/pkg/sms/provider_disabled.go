package sms

import (
	"context"
	"fmt"
)

// DisabledProvider 短信未配置/未启用时的默认实现
type DisabledProvider struct{}

func (p *DisabledProvider) SendCode(_ context.Context, _ string, _ string) error {
	return fmt.Errorf("短信服务未配置")
}
