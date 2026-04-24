package repository

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	"auth-project/internal/domain"
)

// CachedWorklogRepository wraps a WorklogRepository with caching
type CachedWorklogRepository struct {
	repo     domain.WorklogRepository
	cache    map[string]*cacheEntry
	mu       sync.RWMutex
	cacheTTL time.Duration
	cleanup  *time.Ticker
}

type cacheEntry struct {
	data      interface{}
	expiresAt time.Time
}

// NewCachedWorklogRepository creates a new cached worklog repository
func NewCachedWorklogRepository(repo domain.WorklogRepository, cacheTTL time.Duration) *CachedWorklogRepository {
	cachedRepo := &CachedWorklogRepository{
		repo:     repo,
		cache:    make(map[string]*cacheEntry),
		cacheTTL: cacheTTL,
		cleanup:  time.NewTicker(1 * time.Minute),
	}

	// Start cleanup goroutine
	go cachedRepo.cleanupExpired()

	return cachedRepo
}

// Create implements domain.WorklogRepository
func (r *CachedWorklogRepository) Create(worklog *domain.Worklog) error {
	// Invalidate relevant caches on write
	r.invalidateTaskCaches(worklog.TaskID)

	return r.repo.Create(worklog)
}

// FindByTaskID implements domain.WorklogRepository
func (r *CachedWorklogRepository) FindByTaskID(taskID uuid.UUID) ([]*domain.Worklog, error) {
	cacheKey := fmt.Sprintf("worklogs:task:%s", taskID.String())

	// Try cache
	if data, found := r.getFromCache(cacheKey); found {
		if worklogs, ok := data.([]*domain.Worklog); ok {
			return worklogs, nil
		}
	}

	// Fetch from underlying repository
	worklogs, err := r.repo.FindByTaskID(taskID)
	if err != nil {
		return nil, err
	}

	// Store in cache
	r.setInCache(cacheKey, worklogs)

	return worklogs, nil
}

// GetTimeReport implements domain.WorklogRepository
func (r *CachedWorklogRepository) GetTimeReport() ([]*domain.TimeReportRow, error) {
	cacheKey := "report:time"

	// Try cache
	if data, found := r.getFromCache(cacheKey); found {
		if report, ok := data.([]*domain.TimeReportRow); ok {
			return report, nil
		}
	}

	// Fetch from underlying repository
	report, err := r.repo.GetTimeReport()
	if err != nil {
		return nil, err
	}

	// Store in cache (with longer TTL for reports)
	r.mu.Lock()
	r.cache[cacheKey] = &cacheEntry{
		data:      report,
		expiresAt: time.Now().Add(r.cacheTTL * 2), // Reports can be cached longer
	}
	r.mu.Unlock()

	return report, nil
}

// getFromCache retrieves data from cache if not expired
func (r *CachedWorklogRepository) getFromCache(key string) (interface{}, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	entry, exists := r.cache[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(entry.expiresAt) {
		return nil, false
	}

	return entry.data, true
}

// setInCache stores data in cache
func (r *CachedWorklogRepository) setInCache(key string, data interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.cache[key] = &cacheEntry{
		data:      data,
		expiresAt: time.Now().Add(r.cacheTTL),
	}
}

// invalidateTaskCaches removes cached data for a task
func (r *CachedWorklogRepository) invalidateTaskCaches(taskID uuid.UUID) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Invalidate worklogs for this task
	delete(r.cache, fmt.Sprintf("worklogs:task:%s", taskID.String()))

	// Invalidate time report (since it includes worklog data)
	delete(r.cache, "report:time")
}

// cleanupExpired removes expired cache entries
func (r *CachedWorklogRepository) cleanupExpired() {
	for range r.cleanup.C {
		r.mu.Lock()
		now := time.Now()
		for key, entry := range r.cache {
			if now.After(entry.expiresAt) {
				delete(r.cache, key)
			}
		}
		r.mu.Unlock()
	}
}

// Close stops the cleanup ticker
func (r *CachedWorklogRepository) Close() {
	r.cleanup.Stop()
}

// Stats returns cache statistics
func (r *CachedWorklogRepository) Stats() map[string]interface{} {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return map[string]interface{}{
		"cache_size": len(r.cache),
		"cache_ttl":  r.cacheTTL.String(),
	}
}
