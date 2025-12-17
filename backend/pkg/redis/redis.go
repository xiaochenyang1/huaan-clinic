package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"huaan-medical/pkg/config"
)

var client *redis.Client
var enabled bool // Redis是否可用

// Init 初始化Redis连接
func Init(cfg *config.RedisConfig) error {
	if !cfg.Enabled {
		enabled = false
		return nil
	}

	client = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		enabled = false
		return fmt.Errorf("Redis连接失败: %w", err)
	}

	enabled = true
	return nil
}

// TryInit 尝试初始化Redis（失败不报错）
func TryInit(cfg *config.RedisConfig) {
	if !cfg.Enabled {
		enabled = false
		return
	}

	client = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		enabled = false
		client = nil
		return
	}

	enabled = true
}

// IsEnabled 检查Redis是否可用
func IsEnabled() bool {
	return enabled && client != nil
}

// GetClient 获取Redis客户端
func GetClient() *redis.Client {
	return client
}

// Close 关闭Redis连接
func Close() error {
	if client != nil {
		return client.Close()
	}
	return nil
}

// Key 常量定义
const (
	// Token相关
	KeyTokenBlacklist = "token:blacklist:%s" // Token黑名单
	KeyRefreshToken   = "token:refresh:%d"   // 刷新Token

	// 用户相关
	KeyUserInfo = "user:info:%d" // 用户信息缓存

	// 预约相关
	KeyIdempotent     = "idempotent:%s"       // 幂等性Token
	KeyScheduleLock   = "schedule:lock:%d:%s" // 排班锁定
	KeyDailyApptCount = "appt:daily:%d:%s"    // 每日预约计数

	// 验证码相关
	KeySmsCode = "sms:code:%s" // 短信验证码

	// 限流相关
	KeyRateLimit = "rate:limit:%s" // 限流计数
)

// Set 设置键值对
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return client.Set(ctx, key, value, expiration).Err()
}

// Get 获取值
func Get(ctx context.Context, key string) (string, error) {
	return client.Get(ctx, key).Result()
}

// Del 删除键
func Del(ctx context.Context, keys ...string) error {
	return client.Del(ctx, keys...).Err()
}

// Exists 判断键是否存在
func Exists(ctx context.Context, key string) (bool, error) {
	result, err := client.Exists(ctx, key).Result()
	return result > 0, err
}

// Expire 设置过期时间
func Expire(ctx context.Context, key string, expiration time.Duration) error {
	return client.Expire(ctx, key, expiration).Err()
}

// Incr 自增
func Incr(ctx context.Context, key string) (int64, error) {
	return client.Incr(ctx, key).Result()
}

// SetNX 不存在时设置
func SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return client.SetNX(ctx, key, value, expiration).Result()
}

// TTL 获取剩余过期时间
func TTL(ctx context.Context, key string) (time.Duration, error) {
	return client.TTL(ctx, key).Result()
}
