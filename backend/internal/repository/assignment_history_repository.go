package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"auth-project/internal/domain"
)

// AssignmentHistoryRepositoryPostgres implements AssignmentHistoryRepository for PostgreSQL
type AssignmentHistoryRepositoryPostgres struct {
	db *pgxpool.Pool
}

// NewAssignmentHistoryRepositoryPostgres creates a new PostgreSQL assignment history repository
func NewAssignmentHistoryRepositoryPostgres(db *pgxpool.Pool) *AssignmentHistoryRepositoryPostgres {
	return &AssignmentHistoryRepositoryPostgres{db: db}
}

// Create inserts a new assignment history record into the database
func (r *AssignmentHistoryRepositoryPostgres) Create(history *domain.AssignmentHistory) error {
	query := `
		INSERT INTO assignment_history (
			id, task_id, old_assignee_id, new_assignee_id, changed_by, created_at
		) VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(
		context.Background(),
		query,
		history.ID,
		history.TaskID,
		history.OldAssigneeID,
		history.NewAssigneeID,
		history.ChangedBy,
		history.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create assignment history: %w", err)
	}

	return nil
}

// FindByTaskID retrieves all assignment history records for a given task, ordered by created_at descending
func (r *AssignmentHistoryRepositoryPostgres) FindByTaskID(taskID uuid.UUID) ([]*domain.AssignmentHistory, error) {
	query := `
		SELECT 
			id, task_id, old_assignee_id, new_assignee_id, changed_by, created_at
		FROM assignment_history
		WHERE task_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(context.Background(), query, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to query assignment history: %w", err)
	}
	defer rows.Close()

	var histories []*domain.AssignmentHistory
	for rows.Next() {
		var h domain.AssignmentHistory
		err := rows.Scan(
			&h.ID,
			&h.TaskID,
			&h.OldAssigneeID,
			&h.NewAssigneeID,
			&h.ChangedBy,
			&h.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan assignment history row: %w", err)
		}
		histories = append(histories, &h)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating assignment history rows: %w", err)
	}

	return histories, nil
}
