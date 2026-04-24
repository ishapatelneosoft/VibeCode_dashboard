package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Config holds database configuration
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewPool creates a new PostgreSQL connection pool
func NewPool(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	// Optimize for low latency and high concurrency
	poolConfig.MaxConns = 100                           // Increased for high concurrency
	poolConfig.MinConns = 10                            // More connections ready
	poolConfig.MaxConnLifetime = 30 * time.Minute       // Shorter lifetime for fresh connections
	poolConfig.MaxConnIdleTime = 5 * time.Minute        // Faster idle connection cleanup
	poolConfig.HealthCheckPeriod = 15 * time.Second     // More frequent health checks
	poolConfig.MaxConnLifetimeJitter = 30 * time.Second // Jitter for connection recycling

	// Connection acquisition timeout
	poolConfig.ConnConfig.ConnectTimeout = 5 * time.Second

	// TCP keepalive is enabled by default in pgx
	// Statement cache for prepared statements (reduces parse overhead)
	// Note: pgx.QueryExecModeCacheStatement is not available in this version
	// Using default exec mode which is optimized

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection with timeout
	testCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := pool.Ping(testCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Database connection pool established: max=%d, min=%d",
		poolConfig.MaxConns, poolConfig.MinConns)

	return pool, nil
}
