package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type tokenBucket struct {
	tokens    float64
	maxTokens float64
	rate      float64
	lastTime  time.Time
}

type RateLimiter struct {
	mu      sync.Mutex
	buckets map[string]*tokenBucket
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{buckets: make(map[string]*tokenBucket)}
}

func (rl *RateLimiter) Middleware(ratePerMin int) gin.HandlerFunc {
	ratePerSec := float64(ratePerMin) / 60.0
	return func(c *gin.Context) {
		ip := c.ClientIP()

		rl.mu.Lock()
		b, ok := rl.buckets[ip]
		if !ok {
			b = &tokenBucket{
				tokens:    float64(ratePerMin),
				maxTokens: float64(ratePerMin),
				rate:      ratePerSec,
				lastTime:  time.Now(),
			}
			rl.buckets[ip] = b
		}
		rl.mu.Unlock()

		if !b.allow() {
			c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", ratePerMin))
			c.Header("X-RateLimit-Remaining", "0")
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    42901,
				"message": "rate limit exceeded",
				"data":    nil,
			})
			c.Abort()
			return
		}

		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", ratePerMin))
		c.Next()
	}
}

func (b *tokenBucket) allow() bool {
	now := time.Now()
	elapsed := now.Sub(b.lastTime).Seconds()
	b.lastTime = now
	b.tokens += elapsed * b.rate
	if b.tokens > b.maxTokens {
		b.tokens = b.maxTokens
	}
	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}
