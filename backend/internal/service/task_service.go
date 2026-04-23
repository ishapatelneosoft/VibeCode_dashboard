package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"auth-project/internal/domain"
)

// TaskService handles task creation and management business logic
type TaskService struct {
	taskRepo              domain.TaskRepository
	columnRepo            domain.ColumnRepository
	userRepo              domain.UserRepository
	assignmentHistoryRepo domain.AssignmentHistoryRepository
	db                    *pgxpool.Pool
}

// NewTaskService creates a new task service
func NewTaskService(
	taskRepo domain.TaskRepository,
	columnRepo domain.ColumnRepository,
	userRepo domain.UserRepository,
	assignmentHistoryRepo domain.AssignmentHistoryRepository,
	db *pgxpool.Pool,
) *TaskService {
	return &TaskService{
		taskRepo:              taskRepo,
		columnRepo:            columnRepo,
		userRepo:              userRepo,
		assignmentHistoryRepo: assignmentHistoryRepo,
		db:                    db,
	}
}

// CreateTaskRequest represents task creation request
type CreateTaskRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	AssigneeID  *uuid.UUID `json:"assignee_id,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}

// CreateTaskResponse represents task creation response
type CreateTaskResponse struct {
	TaskID   uuid.UUID `json:"task_id"`
	ColumnID uuid.UUID `json:"column_id"`
	Column   string    `json:"column"`
	Position float64   `json:"position"`
	Title    string    `json:"title"`
}

// MoveTaskRequest represents task movement request
type MoveTaskRequest struct {
	DestinationColumnID uuid.UUID  `json:"destination_column_id"`
	BeforeTaskID        *uuid.UUID `json:"before_task_id,omitempty"`
	AfterTaskID         *uuid.UUID `json:"after_task_id,omitempty"`
}

// UpdateAssigneeRequest represents assignee update request
type UpdateAssigneeRequest struct {
	AssigneeID *uuid.UUID `json:"assignee_id,omitempty"`
}

// MoveTaskResponse represents task movement response
type MoveTaskResponse struct {
	TaskID   uuid.UUID `json:"task_id"`
	ColumnID uuid.UUID `json:"column_id"`
	Position float64   `json:"position"`
}

// CreateTask creates a new task with validation and default values
func (s *TaskService) CreateTask(req CreateTaskRequest, createdBy uuid.UUID) (*CreateTaskResponse, error) {
	// Validate title
	if err := s.validateTitle(req.Title); err != nil {
		return nil, err
	}

	// Start transaction for atomic position calculation
	tx, err := s.db.Begin(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(context.Background())

	// Get Backlog column
	backlogColumn, err := s.columnRepo.FindByName("Backlog")
	if err != nil {
		return nil, fmt.Errorf("failed to find backlog column: %w", err)
	}
	if backlogColumn == nil {
		return nil, errors.New("backlog column not found")
	}

	// Get max position in backlog column
	maxPosition, err := s.taskRepo.GetMaxPosition(backlogColumn.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get max position: %w", err)
	}

	// Calculate new position
	newPosition := maxPosition + 1

	// Create task
	task := domain.NewTask(req.Title, backlogColumn.ID, createdBy, newPosition)
	task.Description = req.Description
	task.AssigneeID = req.AssigneeID
	task.DueDate = req.DueDate

	// Save task
	if err := s.taskRepo.Create(task); err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &CreateTaskResponse{
		TaskID:   task.ID,
		ColumnID: backlogColumn.ID,
		Column:   backlogColumn.Name,
		Position: task.Position,
		Title:    task.Title,
	}, nil
}

// GetTask retrieves a task by ID
func (s *TaskService) GetTask(taskID uuid.UUID) (*domain.Task, error) {
	task, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	if task == nil {
		return nil, ErrTaskNotFound
	}
	return task, nil
}

// GetTasksByColumn retrieves all tasks for a column
func (s *TaskService) GetTasksByColumn(columnID uuid.UUID) ([]*domain.Task, error) {
	tasks, err := s.taskRepo.FindByColumn(columnID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	return tasks, nil
}

// UpdateTask updates an existing task
func (s *TaskService) UpdateTask(taskID uuid.UUID, updates map[string]interface{}) error {
	// Get existing task
	task, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}
	if task == nil {
		return ErrTaskNotFound
	}

	// Apply updates
	if title, ok := updates["title"].(string); ok {
		if err := s.validateTitle(title); err != nil {
			return err
		}
		task.Title = title
	}

	if description, ok := updates["description"].(string); ok {
		task.Description = description
	}

	if columnID, ok := updates["column_id"].(uuid.UUID); ok {
		task.ColumnID = columnID
	}

	if assigneeID, ok := updates["assignee_id"].(*uuid.UUID); ok {
		task.AssigneeID = assigneeID
	}

	if dueDate, ok := updates["due_date"].(*time.Time); ok {
		task.DueDate = dueDate
	}

	if position, ok := updates["position"].(float64); ok {
		task.Position = position
	}

	task.UpdatedAt = time.Now()

	// Save updates
	if err := s.taskRepo.Update(task); err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

// MoveTask moves a task to a new column and/or position using fractional indexing
func (s *TaskService) MoveTask(taskID uuid.UUID, req MoveTaskRequest) (*MoveTaskResponse, error) {
	// Start transaction
	tx, err := s.db.Begin(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(context.Background())

	// Fetch task
	task, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to find task: %w", err)
	}
	if task == nil {
		return nil, ErrTaskNotFound
	}

	// Validate destination column
	destColumn, err := s.columnRepo.FindByID(req.DestinationColumnID)
	if err != nil {
		return nil, fmt.Errorf("failed to find destination column: %w", err)
	}
	if destColumn == nil {
		return nil, ErrColumnNotFound
	}

	var beforeTask, afterTask *domain.Task
	// Fetch before task if provided
	if req.BeforeTaskID != nil {
		beforeTask, err = s.taskRepo.FindByID(*req.BeforeTaskID)
		if err != nil {
			return nil, fmt.Errorf("failed to find before task: %w", err)
		}
		if beforeTask == nil {
			return nil, errors.New("before task not found")
		}
		if beforeTask.ColumnID != req.DestinationColumnID {
			return nil, errors.New("before task is not in destination column")
		}
	}
	// Fetch after task if provided
	if req.AfterTaskID != nil {
		afterTask, err = s.taskRepo.FindByID(*req.AfterTaskID)
		if err != nil {
			return nil, fmt.Errorf("failed to find after task: %w", err)
		}
		if afterTask == nil {
			return nil, errors.New("after task not found")
		}
		if afterTask.ColumnID != req.DestinationColumnID {
			return nil, errors.New("after task is not in destination column")
		}
	}

	// Compute new position using fractional indexing
	var newPosition float64
	switch {
	case beforeTask != nil && afterTask != nil:
		// Insert between before and after
		newPosition = (beforeTask.Position + afterTask.Position) / 2
	case beforeTask != nil && afterTask == nil:
		// Insert after before (bottom of column)
		newPosition = beforeTask.Position + 1
	case beforeTask == nil && afterTask != nil:
		// Insert before after (top of column)
		newPosition = afterTask.Position / 2
	default:
		// No before/after provided, move to empty column
		// Use position 0 as default
		newPosition = 0
	}

	// Update task column and position
	task.ColumnID = req.DestinationColumnID
	task.Position = newPosition
	task.UpdatedAt = time.Now()

	// Save changes
	if err := s.taskRepo.Update(task); err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &MoveTaskResponse{
		TaskID:   task.ID,
		ColumnID: task.ColumnID,
		Position: task.Position,
	}, nil
}

// UpdateAssignee updates the assignee of a task and records history
func (s *TaskService) UpdateAssignee(taskID uuid.UUID, newAssigneeID *uuid.UUID, changedBy uuid.UUID) error {
	// Start transaction
	tx, err := s.db.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(context.Background())

	// Fetch task
	task, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		return fmt.Errorf("failed to find task: %w", err)
	}
	if task == nil {
		return ErrTaskNotFound
	}

	// Validate assignee exists if provided
	if newAssigneeID != nil {
		assignee, err := s.userRepo.FindByID(*newAssigneeID)
		if err != nil {
			return fmt.Errorf("failed to find assignee: %w", err)
		}
		if assignee == nil {
			return ErrAssigneeNotFound
		}
	}

	// Check if assignee actually changed
	if (task.AssigneeID == nil && newAssigneeID == nil) ||
		(task.AssigneeID != nil && newAssigneeID != nil && *task.AssigneeID == *newAssigneeID) {
		// No change, skip history
		return nil
	}

	// Record history
	history := domain.NewAssignmentHistory(taskID, changedBy, task.AssigneeID, newAssigneeID)
	if err := s.assignmentHistoryRepo.Create(history); err != nil {
		return fmt.Errorf("failed to create assignment history: %w", err)
	}

	// Update task
	task.AssigneeID = newAssigneeID
	task.UpdatedAt = time.Now()
	if err := s.taskRepo.Update(task); err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(context.Background()); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetAssignmentHistory retrieves assignment history for a task
func (s *TaskService) GetAssignmentHistory(taskID uuid.UUID) ([]*domain.AssignmentHistory, error) {
	// Verify task exists
	task, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to find task: %w", err)
	}
	if task == nil {
		return nil, ErrTaskNotFound
	}

	histories, err := s.assignmentHistoryRepo.FindByTaskID(taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch assignment history: %w", err)
	}
	return histories, nil
}

// DeleteTask removes a task
func (s *TaskService) DeleteTask(taskID uuid.UUID) error {
	if err := s.taskRepo.Delete(taskID); err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	return nil
}

// validateTitle validates task title length
func (s *TaskService) validateTitle(title string) error {
	if len(title) == 0 || len(title) > 255 {
		return ErrInvalidTitle
	}
	return nil
}
