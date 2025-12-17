package middleware

import (
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"huaan-medical/pkg/config"
	"huaan-medical/pkg/errorcode"
	"huaan-medical/pkg/response"
)

// IPLimiter IP限流器
type IPLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

var ipLimiter *IPLimiter

// InitRateLimiter 初始化限流器
func InitRateLimiter(cfg *config.RateLimitConfig) {
	ipLimiter = &IPLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     rate.Limit(cfg.RequestsPerSecond),
		burst:    cfg.Burst,
	}
}

// getLimiter 获取IP对应的限流器
func (l *IPLimiter) getLimiter(ip string) *rate.Limiter {
	l.mu.RLock()
	limiter, exists := l.limiters[ip]
	l.mu.RUnlock()

	if exists {
		return limiter
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// 双重检查
	if limiter, exists = l.limiters[ip]; exists {
		return limiter
	}

	limiter = rate.NewLimiter(l.rate, l.burst)
	l.limiters[ip] = limiter

	return limiter
}

// RateLimit 限流中间件
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := config.Get()
		if cfg == nil || !cfg.RateLimit.Enabled {
			c.Next()
			return
		}

		if ipLimiter == nil {
			c.Next()
			return
		}

		ip := c.ClientIP()
		limiter := ipLimiter.getLimiter(ip)

		if !limiter.Allow() {
			response.Fail(c, errorcode.ErrFrequentRequest)
			c.Abort()
			return
		}

		c.Next()
	}
}

// CleanupLimiters 清理过期的限流器（可定期调用）
func CleanupLimiters() {
	if ipLimiter == nil {
		return
	}

	ipLimiter.mu.Lock()
	defer ipLimiter.mu.Unlock()

	// 简单粗暴地清空，让限流器重新创建
	// 生产环境可以实现更精细的清理策略
	ipLimiter.limiters = make(map[string]*rate.Limiter)
}
