package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"auth-project/internal/domain"
)

// TaskRepositoryPostgres implements TaskRepository for PostgreSQL
type TaskRepositoryPostgres struct {
	db *pgxpool.Pool
}

// NewTaskRepositoryPostgres creates a new PostgreSQL task repository
func NewTaskRepositoryPostgres(db *pgxpool.Pool) *TaskRepositoryPostgres {
	return &TaskRepositoryPostgres{db: db}
}

// Create inserts a new task into the database
func (r *TaskRepositoryPostgres) Create(task *domain.Task) error {
	query := `
		INSERT INTO tasks (
			id, title, description, column_id, assignee_id,
			due_date, created_by, position, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.Exec(
		context.Background(),
		query,
		task.ID,
		task.Title,
		task.Description,
		task.ColumnID,
		task.AssigneeID,
		task.DueDate,
		task.CreatedBy,
		task.Position,
		task.CreatedAt,
		task.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	return nil
}

// FindByID retrieves a task by its ID
func (r *TaskRepositoryPostgres) FindByID(id uuid.UUID) (*domain.Task, error) {
	query := `
		SELECT
			id, title, description, column_id, assignee_id,
			due_date, created_by, position, created_at, updated_at
		FROM tasks
		WHERE id = $1
	`

	var task domain.Task
	var assigneeID *uuid.UUID
	var dueDate *time.Time

	err := r.db.QueryRow(
		context.Background(),
		query,
		id,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.ColumnID,
		&assigneeID,
		&dueDate,
		&task.CreatedBy,
		&task.Position,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find task by ID: %w", err)
	}

	// Convert assigneeID
	if assigneeID != nil {
		task.AssigneeID = assigneeID
	}

	// Convert dueDate
	task.DueDate = dueDate

	return &task, nil
}

// FindByColumn retrieves all tasks for a column ordered by position
func (r *TaskRepositoryPostgres) FindByColumn(columnID uuid.UUID) ([]*domain.Task, error) {
	query := `
		SELECT 
			id, title, description, column_id, assignee_id,
			due_date, created_by, position, created_at, updated_at
		FROM tasks
		WHERE column_id = $1
		ORDER BY position ASC
	`

	rows, err := r.db.Query(context.Background(), query, columnID)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	return r.scanTasks(rows)
}

// FindAllWithColumns retrieves all tasks with their column information
func (r *TaskRepositoryPostgres) FindAllWithColumns() ([]*domain.Task, error) {
	query := `
		SELECT 
			t.id, t.title, t.description, t.column_id, t.assignee_id,
			t.due_date, t.created_by, t.position, t.created_at, t.updated_at
		FROM tasks t
		ORDER BY t.column_id, t.position ASC
	`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to query all tasks: %w", err)
	}
	defer rows.Close()

	return r.scanTasks(rows)
}

// Update updates an existing task
func (r *TaskRepositoryPostgres) Update(task *domain.Task) error {
	query := `
		UPDATE tasks
		SET 
			title = $2,
			description = $3,
			column_id = $4,
			assignee_id = $5,
			due_date = $6,
			position = $7,
			updated_at = $8
		WHERE id = $1
	`

	result, err := r.db.Exec(
		context.Background(),
		query,
		task.ID,
		task.Title,
		task.Description,
		task.ColumnID,
		task.AssigneeID,
		task.DueDate,
		task.Position,
		task.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("task not found")
	}

	return nil
}

// Delete removes a task by ID
func (r *TaskRepositoryPostgres) Delete(id uuid.UUID) error {
	query := `DELETE FROM tasks WHERE id = $1`

	result, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("task not found")
	}

	return nil
}

// GetMaxPosition returns the maximum position value for tasks in a column
func (r *TaskRepositoryPostgres) GetMaxPosition(columnID uuid.UUID) (float64, error) {
	query := `
		SELECT COALESCE(MAX(position), -1.0)
		FROM tasks
		WHERE column_id = $1
	`

	var maxPosition float64
	err := r.db.QueryRow(context.Background(), query, columnID).Scan(&maxPosition)
	if err != nil {
		return -1.0, fmt.Errorf("failed to get max position: %w", err)
	}

	return maxPosition, nil
}

// FindTaskBefore returns the task with the highest position less than the given position in the same column
func (r *TaskRepositoryPostgres) FindTaskBefore(columnID uuid.UUID, position float64) (*domain.Task, error) {
	query := `
		SELECT
			id, title, description, column_id, assignee_id,
			due_date, created_by, position, created_at, updated_at
		FROM tasks
		WHERE column_id = $1 AND position < $2
		ORDER BY position DESC
		LIMIT 1
	`

	rows, err := r.db.Query(context.Background(), query, columnID, position)
	if err != nil {
		return nil, fmt.Errorf("failed to query task before: %w", err)
	}
	defer rows.Close()

	tasks, err := r.scanTasks(rows)
	if err != nil {
		return nil, err
	}
	if len(tasks) == 0 {
		return nil, nil
	}
	return tasks[0], nil
}

// FindTaskAfter returns the task with the lowest position greater than the given position in the same column
func (r *TaskRepositoryPostgres) FindTaskAfter(columnID uuid.UUID, position float64) (*domain.Task, error) {
	query := `
		SELECT
			id, title, description, column_id, assignee_id,
			due_date, created_by, position, created_at, updated_at
		FROM tasks
		WHERE column_id = $1 AND position > $2
		ORDER BY position ASC
		LIMIT 1
	`

	rows, err := r.db.Query(context.Background(), query, columnID, position)
	if err != nil {
		return nil, fmt.Errorf("failed to query task after: %w", err)
	}
	defer rows.Close()

	tasks, err := r.scanTasks(rows)
	if err != nil {
		return nil, err
	}
	if len(tasks) == 0 {
		return nil, nil
	}
	return tasks[0], nil
}

// Helper method to scan task rows
func (r *TaskRepositoryPostgres) scanTasks(rows pgx.Rows) ([]*domain.Task, error) {
	var tasks []*domain.Task

	for rows.Next() {
		var task domain.Task
		var assigneeID *uuid.UUID
		var dueDate *time.Time

		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.ColumnID,
			&assigneeID,
			&dueDate,
			&task.CreatedBy,
			&task.Position,
			&task.CreatedAt,
			&task.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}

		// Convert assigneeID
		if assigneeID != nil {
			task.AssigneeID = assigneeID
		}

		// Convert dueDate
		task.DueDate = dueDate

		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tasks: %w", err)
	}

	return tasks, nil
}
