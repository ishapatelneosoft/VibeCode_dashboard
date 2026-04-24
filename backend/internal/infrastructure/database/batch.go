package database

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// BatchQuery executes multiple queries in parallel with connection pooling
type BatchQuery struct {
	pool *pgxpool.Pool
}

// NewBatchQuery creates a new batch query executor
func NewBatchQuery(pool *pgxpool.Pool) *BatchQuery {
	return &BatchQuery{pool: pool}
}

// BatchResult represents the result of a batched operation
type BatchResult struct {
	Data  interface{}
	Error error
	Index int
}

// ExecuteParallel executes multiple queries in parallel
func (bq *BatchQuery) ExecuteParallel(ctx context.Context, queries []string, args [][]interface{}) ([]BatchResult, error) {
	if len(queries) != len(args) {
		return nil, fmt.Errorf("queries and args length mismatch")
	}

	results := make([]BatchResult, len(queries))
	var wg sync.WaitGroup
	wg.Add(len(queries))

	for i := range queries {
		go func(idx int) {
			defer wg.Done()

			start := time.Now()
			data, err := bq.executeSingle(ctx, queries[idx], args[idx]...)
			elapsed := time.Since(start)

			results[idx] = BatchResult{
				Data:  data,
				Error: err,
				Index: idx,
			}

			// Log slow queries
			if elapsed > 100*time.Millisecond {
				// log.Printf("[SLOW BATCH QUERY] index=%d, query=%s, duration=%v", idx, queries[idx], elapsed)
			}
		}(i)
	}

	wg.Wait()
	return results, nil
}

// executeSingle executes a single query
func (bq *BatchQuery) executeSingle(ctx context.Context, query string, args ...interface{}) (interface{}, error) {
	// Add timeout to prevent hanging queries
	queryCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	rows, err := bq.pool.Query(queryCtx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// For SELECT queries, collect results
	var results []map[string]interface{}
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, err
		}

		fieldDescriptions := rows.FieldDescriptions()
		result := make(map[string]interface{})
		for i, fd := range fieldDescriptions {
			result[string(fd.Name)] = values[i]
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Return single result for single row queries
	if len(results) == 1 {
		return results[0], nil
	}

	return results, nil
}

// BatchLoadTasks loads multiple tasks by ID in a single query
func (bq *BatchQuery) BatchLoadTasks(ctx context.Context, taskIDs []uuid.UUID) ([]map[string]interface{}, error) {
	if len(taskIDs) == 0 {
		return []map[string]interface{}{}, nil
	}

	// Build query with IN clause
	query := `
		SELECT 
			id, title, description, column_id, assignee_id,
			due_date, created_by, position, created_at, updated_at
		FROM tasks 
		WHERE id = ANY($1)
		ORDER BY created_at DESC
	`

	// Convert UUIDs to string array for PostgreSQL
	idStrings := make([]string, len(taskIDs))
	for i, id := range taskIDs {
		idStrings[i] = id.String()
	}

	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := bq.pool.Query(queryCtx, query, idStrings)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []map[string]interface{}
	for rows.Next() {
		var (
			id          uuid.UUID
			title       string
			description *string
			columnID    uuid.UUID
			assigneeID  *uuid.UUID
			dueDate     *time.Time
			createdBy   uuid.UUID
			position    float64
			createdAt   time.Time
			updatedAt   time.Time
		)

		err := rows.Scan(
			&id, &title, &description, &columnID, &assigneeID,
			&dueDate, &createdBy, &position, &createdAt, &updatedAt,
		)
		if err != nil {
			return nil, err
		}

		task := map[string]interface{}{
			"id":          id,
			"title":       title,
			"description": description,
			"column_id":   columnID,
			"assignee_id": assigneeID,
			"due_date":    dueDate,
			"created_by":  createdBy,
			"position":    position,
			"created_at":  createdAt,
			"updated_at":  updatedAt,
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// BatchLoadUsers loads multiple users by ID in a single query
func (bq *BatchQuery) BatchLoadUsers(ctx context.Context, userIDs []uuid.UUID) ([]map[string]interface{}, error) {
	if len(userIDs) == 0 {
		return []map[string]interface{}{}, nil
	}

	query := `
		SELECT id, email, name, created_at
		FROM users
		WHERE id = ANY($1)
	`

	idStrings := make([]string, len(userIDs))
	for i, id := range userIDs {
		idStrings[i] = id.String()
	}

	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := bq.pool.Query(queryCtx, query, idStrings)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var (
			id        uuid.UUID
			email     string
			name      *string
			createdAt time.Time
		)

		err := rows.Scan(&id, &email, &name, &createdAt)
		if err != nil {
			return nil, err
		}

		user := map[string]interface{}{
			"id":         id,
			"email":      email,
			"name":       name,
			"created_at": createdAt,
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// BatchExecutor provides a simple interface for batch operations
type BatchExecutor struct {
	pool *pgxpool.Pool
}

// NewBatchExecutor creates a new batch executor
func NewBatchExecutor(pool *pgxpool.Pool) *BatchExecutor {
	return &BatchExecutor{pool: pool}
}

// ExecuteInTransaction executes multiple operations in a single transaction
func (be *BatchExecutor) ExecuteInTransaction(ctx context.Context, operations []func(tx pgx.Tx) error) error {
	tx, err := be.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, op := range operations {
		if err := op(tx); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
