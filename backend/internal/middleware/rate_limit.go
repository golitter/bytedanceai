package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// ipBucket tracks request count for a single IP within a time window.
type ipBucket struct {
	count   int
	expires time.Time
}

// IPRateLimiter is a simple per-IP sliding window rate limiter.
type IPRateLimiter struct {
	mu      sync.Mutex
	buckets map[string]*ipBucket
	rate    int
	window  time.Duration
}

// NewIPRateLimiter creates a rate limiter allowing at most `rate` requests per `window` per IP.
func NewIPRateLimiter(rate int, window time.Duration) *IPRateLimiter {
	rl := &IPRateLimiter{
		buckets: make(map[string]*ipBucket),
		rate:    rate,
		window:  window,
	}
	go rl.cleanup()
	return rl
}

func (rl *IPRateLimiter) cleanup() {
	ticker := time.NewTicker(rl.window)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, b := range rl.buckets {
			if now.After(b.expires) {
				delete(rl.buckets, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// Middleware returns a Gin middleware that enforces the rate limit.
func (rl *IPRateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		rl.mu.Lock()
		b, ok := rl.buckets[ip]
		now := time.Now()
		if !ok || now.After(b.expires) {
			b = &ipBucket{count: 0, expires: now.Add(rl.window)}
			rl.buckets[ip] = b
		}
		b.count++
		allowed := b.count <= rl.rate
		rl.mu.Unlock()

		if !allowed {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code": http.StatusTooManyRequests,
				"msg":  "too many requests",
			})
			return
		}
		c.Next()
	}
}
