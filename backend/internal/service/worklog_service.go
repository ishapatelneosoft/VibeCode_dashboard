package service

import (
	"errors"
	"fmt"

	"github.com/google/uuid"

	"auth-project/internal/domain"
)

var (
	// ErrInvalidTimeSpent is returned when time_spent is not positive
	ErrInvalidTimeSpent = errors.New("time_spent must be greater than 0")
)

// WorklogService handles worklog creation and retrieval business logic
type WorklogService struct {
	worklogRepo domain.WorklogRepository
	taskRepo    domain.TaskRepository
}

// NewWorklogService creates a new worklog service
func NewWorklogService(
	worklogRepo domain.WorklogRepository,
	taskRepo domain.TaskRepository,
) *WorklogService {
	return &WorklogService{
		worklogRepo: worklogRepo,
		taskRepo:    taskRepo,
	}
}

// LogWorkRequest represents worklog creation request
type LogWorkRequest struct {
	TimeSpent   float64 `json:"time_spent"`
	Description string  `json:"description,omitempty"`
}

// LogWorkResponse represents worklog creation response
type LogWorkResponse struct {
	WorklogID uuid.UUID `json:"worklog_id"`
	TaskID    uuid.UUID `json:"task_id"`
	UserID    uuid.UUID `json:"user_id"`
	TimeSpent float64   `json:"time_spent"`
	CreatedAt string    `json:"created_at"`
}

// LogWork logs time spent on a task by a user
func (s *WorklogService) LogWork(taskID, userID uuid.UUID, req LogWorkRequest) (*LogWorkResponse, error) {
	// Validate time_spent
	if req.TimeSpent <= 0 {
		return nil, ErrInvalidTimeSpent
	}

	// Verify task exists
	task, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to find task: %w", err)
	}
	if task == nil {
		return nil, ErrTaskNotFound
	}

	// Create worklog
	worklog := domain.NewWorklog(taskID, userID, req.TimeSpent, req.Description)

	// Persist
	if err := s.worklogRepo.Create(worklog); err != nil {
		return nil, fmt.Errorf("failed to create worklog: %w", err)
	}

	return &LogWorkResponse{
		WorklogID: worklog.ID,
		TaskID:    worklog.TaskID,
		UserID:    worklog.UserID,
		TimeSpent: worklog.TimeSpent,
		CreatedAt: worklog.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

// GetWorklogs retrieves all worklogs for a task, ordered by created_at descending
func (s *WorklogService) GetWorklogs(taskID uuid.UUID) ([]*domain.Worklog, error) {
	// Verify task exists
	task, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to find task: %w", err)
	}
	if task == nil {
		return nil, ErrTaskNotFound
	}

	worklogs, err := s.worklogRepo.FindByTaskID(taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch worklogs: %w", err)
	}

	return worklogs, nil
}

// TimeReportResponse represents the aggregated time report
type TimeReportResponse struct {
	Tasks      []*domain.TimeReportRow `json:"tasks"`
	GrandTotal float64                 `json:"grand_total"`
}

// GetTimeReport retrieves aggregated time spent per task and grand total
func (s *WorklogService) GetTimeReport() (*TimeReportResponse, error) {
	rows, err := s.worklogRepo.GetTimeReport()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch time report: %w", err)
	}

	// Calculate grand total
	var grandTotal float64
	for _, row := range rows {
		grandTotal += row.TotalHours
	}

	return &TimeReportResponse{
		Tasks:      rows,
		GrandTotal: grandTotal,
	}, nil
}
