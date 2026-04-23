package domain

import (
	"time"

	"github.com/google/uuid"
)

// Task represents a task entity on the shared board
type Task struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	ColumnID    uuid.UUID  `json:"column_id"`
	AssigneeID  *uuid.UUID `json:"assignee_id,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	CreatedBy   uuid.UUID  `json:"created_by"`
	Position    float64    `json:"position"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// NewTask creates a new Task instance with generated ID and timestamps
func NewTask(title string, columnID, createdBy uuid.UUID, position float64) *Task {
	now := time.Now()
	return &Task{
		ID:        uuid.New(),
		Title:     title,
		ColumnID:  columnID,
		CreatedBy: createdBy,
		Position:  position,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// TaskRepository defines the interface for task data operations
type TaskRepository interface {
	Create(task *Task) error
	FindByID(id uuid.UUID) (*Task, error)
	FindByColumn(columnID uuid.UUID) ([]*Task, error)
	FindAllWithColumns() ([]*Task, error)
	Update(task *Task) error
	Delete(id uuid.UUID) error
	GetMaxPosition(columnID uuid.UUID) (float64, error)
	// FindTaskBefore returns the task with the highest position less than the given position in the same column
	FindTaskBefore(columnID uuid.UUID, position float64) (*Task, error)
	// FindTaskAfter returns the task with the lowest position greater than the given position in the same column
	FindTaskAfter(columnID uuid.UUID, position float64) (*Task, error)
}
