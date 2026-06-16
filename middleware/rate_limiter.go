package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	var recent []time.Time
	for _, t := range rl.requests[ip] {
		if t.After(windowStart) {
			recent = append(recent, t)
		}
	}

	if len(recent) >= rl.limit {
		return false
	}

	recent = append(recent, now)
	rl.requests[ip] = recent

	return true
}

func RateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc {
	rl := NewRateLimiter(limit, window)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !rl.Allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
