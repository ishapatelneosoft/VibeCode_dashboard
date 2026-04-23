package domain

import (
	"time"

	"github.com/google/uuid"
)

// Worklog represents a time log entry for a task
type Worklog struct {
	ID          uuid.UUID `json:"id"`
	TaskID      uuid.UUID `json:"task_id"`
	UserID      uuid.UUID `json:"user_id"`
	TimeSpent   float64   `json:"time_spent"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// NewWorklog creates a new Worklog instance
func NewWorklog(taskID, userID uuid.UUID, timeSpent float64, description string) *Worklog {
	return &Worklog{
		ID:          uuid.New(),
		TaskID:      taskID,
		UserID:      userID,
		TimeSpent:   timeSpent,
		Description: description,
		CreatedAt:   time.Now(),
	}
}

// WorklogRepository defines the interface for worklog data operations
type WorklogRepository interface {
	Create(worklog *Worklog) error
	FindByTaskID(taskID uuid.UUID) ([]*Worklog, error)
	GetTimeReport() ([]*TimeReportRow, error)
}
