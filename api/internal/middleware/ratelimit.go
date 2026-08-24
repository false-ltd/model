package middleware

import (
	"fmt"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	maxTrackedIPs   = 10000
	idleSweepAfter  = time.Minute
	bucketIdleLimit = 10 * time.Minute
)

type tokenBucket struct {
	tokens    float64
	maxTokens float64
	rate      float64
	lastTime  time.Time
}

type RateLimiter struct {
	mu        sync.Mutex
	buckets   map[string]*tokenBucket
	lastSweep time.Time
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		buckets:   make(map[string]*tokenBucket),
		lastSweep: time.Now(),
	}
}

func (rl *RateLimiter) Middleware(ratePerMin int) gin.HandlerFunc {
	ratePerSec := float64(ratePerMin) / 60.0
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limit, remaining, reset, ok := rl.allow(ip, ratePerMin, ratePerSec)

		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", reset.Unix()))
		if !ok {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    42901,
				"message": "rate limit exceeded",
				"data":    nil,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// allow consumes one token under the lock so concurrent requests from the
// same IP cannot race on the bucket state, and reports the standard
// X-RateLimit values. Idle buckets are swept periodically to bound memory.
func (rl *RateLimiter) allow(ip string, ratePerMin int, ratePerSec float64) (limit, remaining int, reset time.Time, ok bool) {
	now := time.Now()

	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.sweepLocked(now)

	b, ok := rl.buckets[ip]
	if !ok {
		b = &tokenBucket{
			tokens:    float64(ratePerMin),
			maxTokens: float64(ratePerMin),
			rate:      ratePerSec,
			lastTime:  now,
		}
		rl.buckets[ip] = b
	}

	elapsed := now.Sub(b.lastTime).Seconds()
	b.lastTime = now
	b.tokens += elapsed * b.rate
	if b.tokens > b.maxTokens {
		b.tokens = b.maxTokens
	}

	limit = ratePerMin
	remaining = int(math.Floor(b.tokens))
	refill := (b.maxTokens - b.tokens) / b.rate
	reset = now.Add(time.Duration(refill * float64(time.Second)))

	if b.tokens < 1 {
		return limit, 0, reset, false
	}
	b.tokens--
	remaining = int(math.Floor(b.tokens))
	return limit, remaining, reset, true
}

// sweepLocked evicts buckets that have been idle long enough to have fully
// refilled; must be called with rl.mu held.
func (rl *RateLimiter) sweepLocked(now time.Time) {
	if len(rl.buckets) < maxTrackedIPs || now.Sub(rl.lastSweep) < idleSweepAfter {
		return
	}
	rl.lastSweep = now
	for ip, b := range rl.buckets {
		if b.tokens >= b.maxTokens && now.Sub(b.lastTime) > bucketIdleLimit {
			delete(rl.buckets, ip)
		}
	}
}
