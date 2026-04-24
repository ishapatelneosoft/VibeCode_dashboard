package database

import (
	"context"
	"time"
)

// QueryTimeout provides configurable timeouts for different query types
type QueryTimeout struct {
	// Default timeout for most queries
	Default time.Duration
	// Report queries (aggregations, joins)
	Report time.Duration
	// Simple lookups (by ID, simple filters)
	Lookup time.Duration
	// Write operations (insert, update, delete)
	Write time.Duration
	// Transaction operations
	Transaction time.Duration
}

// DefaultTimeouts returns sensible default timeouts
func DefaultTimeouts() QueryTimeout {
	return QueryTimeout{
		Default:     10 * time.Second,
		Report:      30 * time.Second,
		Lookup:      5 * time.Second,
		Write:       15 * time.Second,
		Transaction: 30 * time.Second,
	}
}

// WithTimeout creates a context with appropriate timeout based on query type
func WithTimeout(ctx context.Context, timeoutType string) (context.Context, context.CancelFunc) {
	timeouts := DefaultTimeouts()

	var timeout time.Duration
	switch timeoutType {
	case "report":
		timeout = timeouts.Report
	case "lookup":
		timeout = timeouts.Lookup
	case "write":
		timeout = timeouts.Write
	case "transaction":
		timeout = timeouts.Transaction
	default:
		timeout = timeouts.Default
	}

	return context.WithTimeout(ctx, timeout)
}

// WithDefaultTimeout creates a context with default timeout
func WithDefaultTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	timeouts := DefaultTimeouts()
	return context.WithTimeout(ctx, timeouts.Default)
}

// ExecuteWithTimeout executes a function with timeout and records timing
func ExecuteWithTimeout(ctx context.Context, timeoutType string, fn func(ctx context.Context) error) error {
	timeoutCtx, cancel := WithTimeout(ctx, timeoutType)
	defer cancel()

	start := time.Now()
	err := fn(timeoutCtx)
	elapsed := time.Since(start)

	// Log slow queries
	if elapsed > 100*time.Millisecond {
		// In production, you would log this to monitoring system
		// log.Printf("[SLOW QUERY] %s took %v", timeoutType, elapsed)
	}

	return err
}

// QueryMetrics tracks query performance
type QueryMetrics struct {
	QueryType    string
	Duration     time.Duration
	Success      bool
	RowsAffected int64
	Error        error
}

// RecordQuery records query metrics for monitoring
func RecordQuery(metrics QueryMetrics) {
	// In production, you would send this to metrics system
	// For now, we just log slow queries
	if metrics.Duration > 500*time.Millisecond {
		status := "success"
		if !metrics.Success {
			status = "error"
		}
		_ = status // Currently unused, but available for future logging
		// log.Printf("[QUERY METRICS] type=%s duration=%v status=%s rows=%d",
		// 	metrics.QueryType, metrics.Duration, status, metrics.RowsAffected)
	}
}
