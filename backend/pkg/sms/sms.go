package sms

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"huaan-medical/pkg/config"
	"huaan-medical/pkg/redis"
	"huaan-medical/pkg/utils"
)

// CodeInfo 验证码信息
type CodeInfo struct {
	Code       string    `json:"code"`
	ExpireAt   time.Time `json:"expire_at"`
	Attempts   int       `json:"attempts"`
	LastSendAt time.Time `json:"last_send_at"`
}

// Service 验证码服务
type Service struct {
	storage sync.Map // 内存存储，降级方案
}

var instance *Service
var once sync.Once

// GetService 获取验证码服务单例
func GetService() *Service {
	once.Do(func() {
		instance = &Service{}
	})
	return instance
}

func (s *Service) getProvider() Provider {
	cfg := config.Get()
	if cfg == nil || !cfg.SMS.Enabled {
		return &DisabledProvider{}
	}
	switch cfg.SMS.Provider {
	case "console":
		return &ConsoleProvider{}
	case "", "disabled":
		return &DisabledProvider{}
	default:
		return &DisabledProvider{}
	}
}

// SendCode 发送验证码
func (s *Service) SendCode(phone string) error {
	// 检查发送频率（60秒内只能发送一次）
	if !s.canSendCode(phone) {
		return fmt.Errorf("验证码发送过于频繁，请稍后再试")
	}

	cfg := config.Get()
	if gin.Mode() == gin.ReleaseMode && cfg != nil && cfg.SMS.Provider == "console" {
		// 避免误把 console provider 用在生产
		return fmt.Errorf("短信服务未配置")
	}

	// 生成6位随机验证码
	code := utils.GenerateRandomNumber(6)

	// 调用短信服务商发送
	if err := s.getProvider().SendCode(context.Background(), phone, code); err != nil {
		return err
	}

	// 保存验证码（发送成功后再保存，避免“发送失败但验证码仍可用”）
	if err := s.saveCode(phone, code); err != nil {
		return err
	}
	return nil
}

// VerifyCode 验证验证码
func (s *Service) VerifyCode(phone, code string) error {
	codeInfo, err := s.getCode(phone)
	if err != nil {
		return fmt.Errorf("验证码不存在或已过期")
	}

	// 检查验证次数（最多3次）
	if codeInfo.Attempts >= 3 {
		s.deleteCode(phone)
		return fmt.Errorf("验证码错误次数过多，请重新获取")
	}

	// 检查是否过期
	if time.Now().After(codeInfo.ExpireAt) {
		s.deleteCode(phone)
		return fmt.Errorf("验证码已过期")
	}

	// 验证码错误
	if codeInfo.Code != code {
		codeInfo.Attempts++
		s.updateCode(phone, codeInfo)
		return fmt.Errorf("验证码错误")
	}

	// 验证成功，删除验证码
	s.deleteCode(phone)
	return nil
}

// canSendCode 检查是否可以发送验证码
func (s *Service) canSendCode(phone string) bool {
	codeInfo, err := s.getCode(phone)
	if err != nil {
		return true // 不存在，可以发送
	}

	// 60秒内不能重复发送
	return time.Now().Sub(codeInfo.LastSendAt) >= 60*time.Second
}

// saveCode 保存验证码
func (s *Service) saveCode(phone, code string) error {
	codeInfo := &CodeInfo{
		Code:       code,
		ExpireAt:   time.Now().Add(5 * time.Minute),
		Attempts:   0,
		LastSendAt: time.Now(),
	}

	// 优先使用Redis
	if redis.IsEnabled() {
		ctx := context.Background()
		key := fmt.Sprintf(redis.KeySmsCode, phone)
		return redis.Set(ctx, key, code, 5*time.Minute)
	}

	// 降级到内存存储
	s.storage.Store(phone, codeInfo)
	return nil
}

// getCode 获取验证码
func (s *Service) getCode(phone string) (*CodeInfo, error) {
	// 优先从Redis获取
	if redis.IsEnabled() {
		ctx := context.Background()
		key := fmt.Sprintf(redis.KeySmsCode, phone)
		code, err := redis.Get(ctx, key)
		if err == nil {
			return &CodeInfo{
				Code:     code,
				ExpireAt: time.Now().Add(5 * time.Minute), // Redis已处理过期
				Attempts: 0,
			}, nil
		}
	}

	// 从内存获取
	value, ok := s.storage.Load(phone)
	if !ok {
		return nil, fmt.Errorf("验证码不存在")
	}

	return value.(*CodeInfo), nil
}

// updateCode 更新验证码信息
func (s *Service) updateCode(phone string, info *CodeInfo) {
	if !redis.IsEnabled() {
		s.storage.Store(phone, info)
	}
	// Redis模式下不需要更新attempts（简化实现）
}

// deleteCode 删除验证码
func (s *Service) deleteCode(phone string) {
	if redis.IsEnabled() {
		ctx := context.Background()
		key := fmt.Sprintf(redis.KeySmsCode, phone)
		redis.Del(ctx, key)
	}
	s.storage.Delete(phone)
}
