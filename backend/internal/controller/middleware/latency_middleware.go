package middleware

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LatencyMiddleware captures request latency and logs slow requests
func LatencyMiddleware(threshold time.Duration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		ctx.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get status code
		status := ctx.Writer.Status()

		// Log based on latency threshold
		if latency > threshold {
			log.Printf("[WARN] SLOW REQUEST: %s %s took %v (status: %d, client: %s)",
				ctx.Request.Method, ctx.Request.URL.Path, latency, status, ctx.ClientIP())
		} else if status >= 400 {
			log.Printf("[WARN] ERROR REQUEST: %s %s status %d (latency: %v)",
				ctx.Request.Method, ctx.Request.URL.Path, status, latency)
		} else {
			log.Printf("[INFO] Request: %s %s completed in %v (status: %d)",
				ctx.Request.Method, ctx.Request.URL.Path, latency, status)
		}

		// Add latency header for client visibility
		ctx.Header("X-Response-Time", fmt.Sprintf("%dms", latency.Milliseconds()))

		// Add server timing header (W3C Server-Timing spec)
		ctx.Header("Server-Timing", fmt.Sprintf("total;dur=%d", latency.Milliseconds()))
	}
}

// QueryLatencyMiddleware captures database query latency
func QueryLatencyMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Store query timings in context
		queryTimings := make([]time.Duration, 0)
		ctx.Set("query_timings", &queryTimings)

		ctx.Next()

		// Log total query time if available
		if timings, exists := ctx.Get("query_timings"); exists {
			if timingSlice, ok := timings.(*[]time.Duration); ok && len(*timingSlice) > 0 {
				var totalQueryTime time.Duration
				for _, t := range *timingSlice {
					totalQueryTime += t
				}

				avgQueryTime := totalQueryTime / time.Duration(len(*timingSlice))
				log.Printf("[DEBUG] Database queries: %s %s - %d queries, total: %v, avg: %v",
					ctx.Request.Method, ctx.Request.URL.Path,
					len(*timingSlice), totalQueryTime, avgQueryTime)

				// Add query timing header
				ctx.Header("X-Query-Time", fmt.Sprintf("%dms", totalQueryTime.Milliseconds()))
			}
		}
	}
}

// RecordQueryTiming records a single query execution time
func RecordQueryTiming(ctx *gin.Context, duration time.Duration) {
	if timings, exists := ctx.Get("query_timings"); exists {
		if timingSlice, ok := timings.(*[]time.Duration); ok {
			*timingSlice = append(*timingSlice, duration)
		}
	}
}

// ResponseSizeMiddleware logs response size for optimization
func ResponseSizeMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Use custom response writer to capture size
		writer := &responseSizeWriter{
			ResponseWriter: ctx.Writer,
			size:           0,
		}
		ctx.Writer = writer

		ctx.Next()

		// Log response size for large responses
		if writer.size > 1024*1024 { // 1MB
			sizeMB := float64(writer.size) / (1024 * 1024)
			log.Printf("[WARN] Large response: %s %s - %.2fMB",
				ctx.Request.Method, ctx.Request.URL.Path, sizeMB)
		}

		// Add response size header
		ctx.Header("X-Response-Size", fmt.Sprintf("%d", writer.size))
	}
}

// responseSizeWriter captures response size
type responseSizeWriter struct {
	gin.ResponseWriter
	size int
}

func (w *responseSizeWriter) Write(b []byte) (int, error) {
	w.size += len(b)
	return w.ResponseWriter.Write(b)
}

func (w *responseSizeWriter) WriteString(s string) (int, error) {
	w.size += len(s)
	return w.ResponseWriter.WriteString(s)
}
