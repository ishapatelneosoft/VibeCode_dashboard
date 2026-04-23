package domain

import (
	"github.com/google/uuid"
)

// TimeReportRow represents a single row in the time report aggregation
type TimeReportRow struct {
	TaskID        uuid.UUID  `json:"task_id"`
	Title         string     `json:"title"`
	ColumnName    string     `json:"column_name"`
	AssigneeID    *uuid.UUID `json:"assignee_id,omitempty"`
	AssigneeEmail *string    `json:"assignee_email,omitempty"`
	TotalHours    float64    `json:"total_hours"`
}
