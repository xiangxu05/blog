package middleware

import (
	"blog/internal/message"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// IPRateLimiter IP限流器
type IPRateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit // 每秒允许的请求数
	burst    int        // 突发请求数
}

// NewIPRateLimiter 创建IP限流器
// rps: 每秒请求数 (requests per second)
// burst: 突发请求数
func NewIPRateLimiter(rps float64, burst int) *IPRateLimiter {
	return &IPRateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     rate.Limit(rps),
		burst:    burst,
	}
}

// getLimiter 获取或创建指定IP的限流器
func (rl *IPRateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.RLock()
	limiter, exists := rl.limiters[ip]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		defer rl.mu.Unlock()
		// 双重检查，避免重复创建
		limiter, exists = rl.limiters[ip]
		if !exists {
			limiter = rate.NewLimiter(rl.rate, rl.burst)
			rl.limiters[ip] = limiter
		}
	}

	return limiter
}

// cleanup 定期清理不活跃的限流器（避免内存泄漏）
func (rl *IPRateLimiter) cleanup() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		// 简单清理：如果map太大就清空（实际可以更智能）
		if len(rl.limiters) > 10000 {
			rl.limiters = make(map[string]*rate.Limiter)
		}
		rl.mu.Unlock()
	}
}

// IPRateLimit IP限流中间件
// rps: 每秒请求数，默认 100
// burst: 突发请求数，默认 200
func IPRateLimit(rps float64, burst int) gin.HandlerFunc {
	if rps <= 0 {
		rps = 100 // 默认每秒100个请求
	}
	if burst <= 0 {
		burst = 200 // 默认突发200个请求
	}

	limiter := NewIPRateLimiter(rps, burst)
	// 启动清理协程
	go limiter.cleanup()

	return func(c *gin.Context) {
		// 获取客户端IP
		ip := c.ClientIP()
		if ip == "" {
			ip = c.RemoteIP()
		}

		// 获取该IP的限流器
		limiter := limiter.getLimiter(ip)

		// 检查是否允许请求
		if !limiter.Allow() {
			message.SendMsg(c, http.StatusTooManyRequests, "请求过于频繁，请稍后再试", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

// LoginRateLimiter 登录接口限流器（更严格）
type LoginRateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

// NewLoginRateLimiter 创建登录限流器
// rps: 每秒请求数，默认 0.1 (每10秒1次)
// burst: 突发请求数，默认 3
func NewLoginRateLimiter(rps float64, burst int) *LoginRateLimiter {
	return &LoginRateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     rate.Limit(rps),
		burst:    burst,
	}
}

// getLimiter 获取或创建指定IP的限流器
func (rl *LoginRateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.RLock()
	limiter, exists := rl.limiters[ip]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		defer rl.mu.Unlock()
		limiter, exists = rl.limiters[ip]
		if !exists {
			limiter = rate.NewLimiter(rl.rate, rl.burst)
			rl.limiters[ip] = limiter
		}
	}

	return limiter
}

// cleanup 定期清理
func (rl *LoginRateLimiter) cleanup() {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		if len(rl.limiters) > 5000 {
			rl.limiters = make(map[string]*rate.Limiter)
		}
		rl.mu.Unlock()
	}
}

// LoginRateLimit 登录接口限流中间件（更严格）
// rps: 每秒请求数，默认 0.1 (每10秒1次)
// burst: 突发请求数，默认 3
func LoginRateLimit(rps float64, burst int) gin.HandlerFunc {
	if rps <= 0 {
		rps = 0.1 // 默认每10秒1次请求
	}
	if burst <= 0 {
		burst = 3 // 默认突发3次
	}

	limiter := NewLoginRateLimiter(rps, burst)
	// 启动清理协程
	go limiter.cleanup()

	return func(c *gin.Context) {
		// 获取客户端IP
		ip := c.ClientIP()
		if ip == "" {
			ip = c.RemoteIP()
		}

		// 获取该IP的限流器
		limiter := limiter.getLimiter(ip)

		// 检查是否允许请求
		if !limiter.Allow() {
			message.SendMsg(c, http.StatusTooManyRequests, "登录请求过于频繁，请稍后再试", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
