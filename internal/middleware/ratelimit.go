package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type tokenBucket struct {
	tokens int
	max    int
	last   time.Time
	rate   time.Duration
	mu     sync.Mutex
}

func (tb *tokenBucket) take() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.last)
	tokensToAdd := int(elapsed / tb.rate)
	if tokensToAdd > 0 {
		tb.tokens = min(tb.tokens+tokensToAdd, tb.max)
		tb.last = now
	}

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

var buckets sync.Map

func getBucket(key string, maxTokens int, rate time.Duration) *tokenBucket {
	if v, ok := buckets.Load(key); ok {
		return v.(*tokenBucket)
	}
	bucket := &tokenBucket{
		tokens: maxTokens,
		max:    maxTokens,
		last:   time.Now(),
		rate:   rate,
	}
	buckets.Store(key, bucket)
	return bucket
}

func RateLimit(maxTokens int, rate time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		bucket := getBucket(ip, maxTokens, rate)
		if !bucket.take() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code": 429,
				"msg":  "请求过于频繁，请稍后再试",
				"data": nil,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
