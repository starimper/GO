package utils

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)


const (
	RateLimitMaxRequests = 5            
	RateLimitWindow      = 1 * time.Minute 
)

type bucket struct {
	count       int
	windowStart time.Time
}

type RateLimiter struct {
	mu      sync.Mutex
	buckets map[string]*bucket
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		buckets: make(map[string]*bucket),
	}
}

func (rl *RateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	b, exists := rl.buckets[key]

	if !exists || now.Sub(b.windowStart) > RateLimitWindow {
		rl.buckets[key] = &bucket{count: 1, windowStart: now}
		return true
	}

	b.count++
	return b.count <= RateLimitMaxRequests
}

var globalLimiter = NewRateLimiter()

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var key string

		if userID, ok := parseTokenFromHeader(c); ok {
			key = "user:" + userID
		} else {
			key = "ip:" + c.ClientIP()
		}

		if !globalLimiter.allow(key) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please wait before retrying.",
				"limit": RateLimitMaxRequests,
				"window": RateLimitWindow.String(),
			})
			return
		}

		c.Next()
	}
}
