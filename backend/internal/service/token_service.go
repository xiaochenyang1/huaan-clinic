package service

import (
	"context"
	"fmt"
	"time"

	"huaan-medical/pkg/errorcode"
	"huaan-medical/pkg/redis"
	"huaan-medical/pkg/utils"
)

// TokenService Token服务
type TokenService struct{}

// NewTokenService 创建Token服务实例
func NewTokenService() *TokenService {
	return &TokenService{}
}

// GenerateIdempotentToken 生成幂等Token
func (s *TokenService) GenerateIdempotentToken(userID int64) (string, int64, error) {
	// 生成随机Token
	token := utils.GenerateShortUUID()

	// Token key: idempotent:用户ID:Token
	key := fmt.Sprintf("idempotent:%d:%s", userID, token)

	// 存储到Redis，有效期5分钟
	expiresIn := int64(300) // 5分钟 = 300秒
	ctx := context.Background()
	rdb := redis.GetClient()

	err := rdb.Set(ctx, key, "1", time.Duration(expiresIn)*time.Second).Err()
	if err != nil {
		return "", 0, errorcode.NewWithMessage(errorcode.ErrInternalServer, "生成Token失败")
	}

	return token, expiresIn, nil
}

// ValidateAndConsumeIdempotentToken 验证并消费幂等Token
func (s *TokenService) ValidateAndConsumeIdempotentToken(userID int64, token string) error {
	if token == "" {
		return errorcode.NewWithMessage(errorcode.ErrInvalidParams, "幂等Token不能为空")
	}

	key := fmt.Sprintf("idempotent:%d:%s", userID, token)
	ctx := context.Background()
	rdb := redis.GetClient()

	// 尝试删除Token（原子操作）
	result, err := rdb.Del(ctx, key).Result()
	if err != nil {
		return errorcode.NewWithMessage(errorcode.ErrInternalServer, "验证Token失败")
	}

	// 如果删除返回0，说明Token不存在或已被消费
	if result == 0 {
		return errorcode.NewWithMessage(errorcode.ErrInvalidParams, "Token无效或已过期")
	}

	return nil
}
