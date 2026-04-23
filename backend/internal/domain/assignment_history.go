package domain

import (
	"time"

	"github.com/google/uuid"
)

// AssignmentHistory represents a record of assignee changes for a task
type AssignmentHistory struct {
	ID            uuid.UUID  `json:"id"`
	TaskID        uuid.UUID  `json:"task_id"`
	OldAssigneeID *uuid.UUID `json:"old_assignee_id,omitempty"`
	NewAssigneeID *uuid.UUID `json:"new_assignee_id,omitempty"`
	ChangedBy     uuid.UUID  `json:"changed_by"`
	CreatedAt     time.Time  `json:"created_at"`
}

// NewAssignmentHistory creates a new AssignmentHistory instance
func NewAssignmentHistory(taskID, changedBy uuid.UUID, oldAssigneeID, newAssigneeID *uuid.UUID) *AssignmentHistory {
	return &AssignmentHistory{
		ID:            uuid.New(),
		TaskID:        taskID,
		OldAssigneeID: oldAssigneeID,
		NewAssigneeID: newAssigneeID,
		ChangedBy:     changedBy,
		CreatedAt:     time.Now(),
	}
}

// AssignmentHistoryRepository defines the interface for assignment history data operations
type AssignmentHistoryRepository interface {
	Create(history *AssignmentHistory) error
	FindByTaskID(taskID uuid.UUID) ([]*AssignmentHistory, error)
}
