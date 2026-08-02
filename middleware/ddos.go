package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// NewGlobalRateLimitMiddleware limits total request throughput.
func NewGlobalRateLimitMiddleware(r rate.Limit, b int) gin.HandlerFunc {
	limiter := rate.NewLimiter(r, b)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests",
			})
			return
		}

		c.Next()
	}
}

// NewMaxBodyBytesMiddleware rejects requests over the specified body size.
func NewMaxBodyBytesMiddleware(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxBytes {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": "request body too large",
			})
			return
		}

		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		c.Next()
	}
}

// NewConcurrencyLimitMiddleware caps in-flight requests.
func NewConcurrencyLimitMiddleware(maxConcurrent int) gin.HandlerFunc {
	sem := make(chan struct{}, maxConcurrent)

	return func(c *gin.Context) {
		select {
		case sem <- struct{}{}:
			defer func() {
				time.Sleep(10 * time.Millisecond)
				<-sem
			}()
		default:
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				"error": "server busy",
			})
			return
		}

		c.Next()
	}
}
