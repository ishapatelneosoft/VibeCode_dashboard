package service

import (
	"errors"
	"fmt"

	"github.com/google/uuid"

	"auth-project/internal/domain"
)

var (
	// ErrInvalidTitle is returned when title validation fails
	ErrInvalidTitle = errors.New("title must be between 1 and 255 characters")
	// ErrColumnNotFound is returned when specified column doesn't exist
	ErrColumnNotFound = errors.New("column not found")
	// ErrTaskNotFound is returned when task doesn't exist
	ErrTaskNotFound = errors.New("task not found")
	// ErrAssigneeNotFound is returned when assignee user doesn't exist
	ErrAssigneeNotFound = errors.New("assignee not found")
)

// BoardService handles board and task business logic
type BoardService struct {
	columnRepo domain.ColumnRepository
	taskRepo   domain.TaskRepository
}

// NewBoardService creates a new board service
func NewBoardService(
	columnRepo domain.ColumnRepository,
	taskRepo domain.TaskRepository,
) *BoardService {
	return &BoardService{
		columnRepo: columnRepo,
		taskRepo:   taskRepo,
	}
}

// BoardColumn represents a column with its tasks for board view
type BoardColumn struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Order int       `json:"order"`
	Tasks []*Task   `json:"tasks"`
}

// Task represents a task DTO for board view
type Task struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	ColumnID    uuid.UUID  `json:"column_id"`
	AssigneeID  *uuid.UUID `json:"assignee_id,omitempty"`
	DueDate     *string    `json:"due_date,omitempty"`
	CreatedBy   uuid.UUID  `json:"created_by"`
	Position    float64    `json:"position"`
	CreatedAt   string     `json:"created_at"`
	UpdatedAt   string     `json:"updated_at"`
}

// BoardResponse represents the complete board response
type BoardResponse struct {
	Columns []*BoardColumn `json:"columns"`
}

// GetBoard retrieves the complete board with all columns and tasks
func (s *BoardService) GetBoard() (*BoardResponse, error) {
	// Get all columns in order
	columns, err := s.columnRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	// Get all tasks with their columns
	tasks, err := s.taskRepo.FindAllWithColumns()
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	// Group tasks by column ID
	tasksByColumn := make(map[uuid.UUID][]*Task)
	for _, task := range tasks {
		taskDTO := s.domainTaskToDTO(task)
		tasksByColumn[task.ColumnID] = append(tasksByColumn[task.ColumnID], taskDTO)
	}

	// Build board columns
	boardColumns := make([]*BoardColumn, 0, len(columns))
	for _, column := range columns {
		boardColumn := &BoardColumn{
			ID:    column.ID,
			Name:  column.Name,
			Order: column.Order,
			Tasks: tasksByColumn[column.ID],
		}
		boardColumns = append(boardColumns, boardColumn)
	}

	return &BoardResponse{
		Columns: boardColumns,
	}, nil
}

// ValidateTitle validates task title length
func (s *BoardService) ValidateTitle(title string) error {
	if len(title) == 0 || len(title) > 255 {
		return ErrInvalidTitle
	}
	return nil
}

// Helper method to convert domain task to DTO
func (s *BoardService) domainTaskToDTO(task *domain.Task) *Task {
	taskDTO := &Task{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		ColumnID:    task.ColumnID,
		AssigneeID:  task.AssigneeID,
		CreatedBy:   task.CreatedBy,
		Position:    task.Position,
		CreatedAt:   task.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   task.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if task.DueDate != nil {
		dueDateStr := task.DueDate.Format("2006-01-02T15:04:05Z07:00")
		taskDTO.DueDate = &dueDateStr
	}

	return taskDTO
}
