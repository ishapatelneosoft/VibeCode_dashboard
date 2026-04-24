package controller

import (
	"context"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// MetricsController handles performance metrics endpoints
type MetricsController struct {
	db *pgxpool.Pool
}

// NewMetricsController creates a new metrics controller
func NewMetricsController(db *pgxpool.Pool) *MetricsController {
	return &MetricsController{db: db}
}

// GetMetrics returns system performance metrics
func (c *MetricsController) GetMetrics(ctx *gin.Context) {
	// Get database pool statistics
	dbStats := c.db.Stat()

	// Get memory statistics
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// Get goroutine count
	goroutineCount := runtime.NumGoroutine()

	// Calculate system uptime (simplified)
	uptime := time.Since(startTime)

	metrics := gin.H{
		"timestamp": time.Now().UTC().Format(time.RFC3339),

		// Database metrics
		"database": gin.H{
			"max_connections":          dbStats.MaxConns(),
			"total_connections":        dbStats.TotalConns(),
			"idle_connections":         dbStats.IdleConns(),
			"acquired_connections":     dbStats.AcquiredConns(),
			"constructing_connections": dbStats.ConstructingConns(),
			"empty_acquires":           dbStats.EmptyAcquireCount(),
			"canceled_acquires":        dbStats.CanceledAcquireCount(),
			"acquire_duration":         dbStats.AcquireDuration().String(),
		},

		// Memory metrics
		"memory": gin.H{
			"alloc_bytes":       memStats.Alloc,
			"total_alloc_bytes": memStats.TotalAlloc,
			"sys_bytes":         memStats.Sys,
			"heap_alloc_bytes":  memStats.HeapAlloc,
			"heap_sys_bytes":    memStats.HeapSys,
			"heap_idle_bytes":   memStats.HeapIdle,
			"heap_inuse_bytes":  memStats.HeapInuse,
			"num_gc":            memStats.NumGC,
			"last_gc":           memStats.LastGC,
			"gc_pause_total_ns": memStats.PauseTotalNs,
		},

		// Goroutine metrics
		"goroutines": goroutineCount,

		// System metrics
		"system": gin.H{
			"num_cpu":        runtime.NumCPU(),
			"gomaxprocs":     runtime.GOMAXPROCS(0),
			"uptime_seconds": int64(uptime.Seconds()),
		},

		// Application metrics (would be populated from middleware)
		"application": gin.H{
			"active_requests": 0, // Would track from middleware
			"total_requests":  0, // Would track from middleware
			"error_rate":      0, // Would track from middleware
		},
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data:    metrics,
	})
}

// GetHealth returns detailed health check with performance indicators
func (c *MetricsController) GetHealth(ctx *gin.Context) {
	healthStatus := gin.H{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"checks":    gin.H{},
	}

	// Check database connectivity
	dbCtx, cancel := context.WithTimeout(ctx.Request.Context(), 2*time.Second)
	defer cancel()

	if err := c.db.Ping(dbCtx); err != nil {
		healthStatus["status"] = "unhealthy"
		healthStatus["checks"].(gin.H)["database"] = gin.H{
			"status": "unhealthy",
			"error":  err.Error(),
		}
	} else {
		healthStatus["checks"].(gin.H)["database"] = gin.H{
			"status":     "healthy",
			"latency_ms": 0, // Would measure actual ping time
		}
	}

	// Check memory pressure
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	memoryUsagePercent := float64(memStats.HeapInuse) / float64(memStats.HeapSys) * 100
	if memoryUsagePercent > 90 {
		healthStatus["status"] = "degraded"
		healthStatus["checks"].(gin.H)["memory"] = gin.H{
			"status":        "degraded",
			"usage_percent": memoryUsagePercent,
			"warning":       "high memory usage",
		}
	} else {
		healthStatus["checks"].(gin.H)["memory"] = gin.H{
			"status":        "healthy",
			"usage_percent": memoryUsagePercent,
		}
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data:    healthStatus,
	})
}

// GetPerformance returns performance benchmarks
func (c *MetricsController) GetPerformance(ctx *gin.Context) {
	// This would track request latency percentiles, error rates, etc.
	// For now, return basic performance indicators

	performance := gin.H{
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"indicators": gin.H{
			"request_latency_p50_ms": 0, // Would track from middleware
			"request_latency_p95_ms": 0, // Would track from middleware
			"request_latency_p99_ms": 0, // Would track from middleware
			"database_query_p50_ms":  0, // Would track from middleware
			"database_query_p95_ms":  0, // Would track from middleware
			"database_query_p99_ms":  0, // Would track from middleware
			"error_rate_5m":          0, // Would track from middleware
			"throughput_rps":         0, // Would track from middleware
		},
		"recommendations": []string{
			"Enable response compression for large payloads",
			"Implement query caching for frequently accessed data",
			"Use connection pooling with appropriate limits",
			"Add database indexes for slow queries",
		},
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data:    performance,
	})
}

// startTime tracks when the application started
var startTime = time.Now()
