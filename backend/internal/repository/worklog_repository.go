package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"auth-project/internal/domain"
)

// WorklogRepositoryPostgres implements WorklogRepository for PostgreSQL
type WorklogRepositoryPostgres struct {
	db *pgxpool.Pool
}

// NewWorklogRepositoryPostgres creates a new PostgreSQL worklog repository
func NewWorklogRepositoryPostgres(db *pgxpool.Pool) *WorklogRepositoryPostgres {
	return &WorklogRepositoryPostgres{db: db}
}

// Create inserts a new worklog record into the database
func (r *WorklogRepositoryPostgres) Create(worklog *domain.Worklog) error {
	query := `
		INSERT INTO worklogs (
			id, task_id, user_id, time_spent, description, created_at
		) VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(
		context.Background(),
		query,
		worklog.ID,
		worklog.TaskID,
		worklog.UserID,
		worklog.TimeSpent,
		worklog.Description,
		worklog.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create worklog: %w", err)
	}

	return nil
}

// FindByTaskID retrieves all worklog records for a given task, ordered by created_at descending
func (r *WorklogRepositoryPostgres) FindByTaskID(taskID uuid.UUID) ([]*domain.Worklog, error) {
	query := `
		SELECT 
			id, task_id, user_id, time_spent, description, created_at
		FROM worklogs
		WHERE task_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(context.Background(), query, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to query worklogs: %w", err)
	}
	defer rows.Close()

	var worklogs []*domain.Worklog
	for rows.Next() {
		var w domain.Worklog
		err := rows.Scan(
			&w.ID,
			&w.TaskID,
			&w.UserID,
			&w.TimeSpent,
			&w.Description,
			&w.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan worklog row: %w", err)
		}
		worklogs = append(worklogs, &w)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating worklog rows: %w", err)
	}

	return worklogs, nil
}

// GetTimeReport retrieves aggregated time spent per task across the project
func (r *WorklogRepositoryPostgres) GetTimeReport() ([]*domain.TimeReportRow, error) {
	query := `
		SELECT
			t.id AS task_id,
			t.title,
			c.name AS column_name,
			t.assignee_id,
			u.email AS assignee_email,
			COALESCE(SUM(w.time_spent), 0) AS total_hours
		FROM tasks t
		JOIN columns c ON t.column_id = c.id
		LEFT JOIN users u ON t.assignee_id = u.id
		LEFT JOIN worklogs w ON t.id = w.task_id
		GROUP BY t.id, t.title, c.name, t.assignee_id, u.email
		ORDER BY c.name, t.title
	`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to query time report: %w", err)
	}
	defer rows.Close()

	var report []*domain.TimeReportRow
	for rows.Next() {
		var row domain.TimeReportRow
		var assigneeID *uuid.UUID
		var assigneeEmail *string
		err := rows.Scan(
			&row.TaskID,
			&row.Title,
			&row.ColumnName,
			&assigneeID,
			&assigneeEmail,
			&row.TotalHours,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan time report row: %w", err)
		}
		row.AssigneeID = assigneeID
		row.AssigneeEmail = assigneeEmail
		report = append(report, &row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating time report rows: %w", err)
	}

	return report, nil
}
