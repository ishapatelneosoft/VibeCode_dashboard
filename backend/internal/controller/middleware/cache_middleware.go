package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// CacheMiddleware sets cache control headers for responses
// Useful for reducing database queries on repeated requests
func CacheMiddleware(duration time.Duration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Set cache control headers
		maxAge := fmt.Sprintf("private, max-age=%d", int(duration.Seconds()))
		ctx.Header("Cache-Control", maxAge)
		ctx.Header("Expires", time.Now().Add(duration).Format(time.RFC1123))

		ctx.Next()
	}
}

// NoCache middleware prevents caching of responses
func NoCache() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		ctx.Header("Pragma", "no-cache")
		ctx.Header("Expires", "0")

		ctx.Next()
	}
}
