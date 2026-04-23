package tests

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"auth-project/internal/domain"
	"auth-project/internal/service"
)

// MockAssignmentHistoryRepository is a mock implementation of AssignmentHistoryRepository
type MockAssignmentHistoryRepository struct {
	mock.Mock
}

func (m *MockAssignmentHistoryRepository) Create(history *domain.AssignmentHistory) error {
	args := m.Called(history)
	return args.Error(0)
}

func (m *MockAssignmentHistoryRepository) FindByTaskID(taskID uuid.UUID) ([]*domain.AssignmentHistory, error) {
	args := m.Called(taskID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.AssignmentHistory), args.Error(1)
}

func TestTaskService_UpdateAssignee_Success_Assign(t *testing.T) {
	// Setup mocks
	taskRepo := new(MockTaskRepository)
	columnRepo := new(MockColumnRepository)
	userRepo := new(MockUserRepository)
	assignmentHistoryRepo := new(MockAssignmentHistoryRepository)
	// db mock not needed for unit test (transaction is mocked via repo)
	// We'll skip db transaction for simplicity
	taskService := service.NewTaskService(taskRepo, columnRepo, userRepo, assignmentHistoryRepo, nil)

	taskID := uuid.New()
	assigneeID := uuid.New()
	changedBy := uuid.New()

	task := &domain.Task{
		ID:         taskID,
		AssigneeID: nil,
		UpdatedAt:  time.Now(),
	}
	assignee := &domain.User{
		ID:    assigneeID,
		Email: "user@example.com",
	}

	// Mock expectations
	taskRepo.On("FindByID", taskID).Return(task, nil)
	userRepo.On("FindByID", assigneeID).Return(assignee, nil)
	assignmentHistoryRepo.On("Create", mock.AnythingOfType("*domain.AssignmentHistory")).Return(nil)
	taskRepo.On("Update", mock.AnythingOfType("*domain.Task")).Return(nil)

	// Execute
	err := taskService.UpdateAssignee(taskID, &assigneeID, changedBy)

	// Assert
	assert.NoError(t, err)
	taskRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
	assignmentHistoryRepo.AssertExpectations(t)
}

func TestTaskService_UpdateAssignee_Success_Unassign(t *testing.T) {
	taskRepo := new(MockTaskRepository)
	columnRepo := new(MockColumnRepository)
	userRepo := new(MockUserRepository)
	assignmentHistoryRepo := new(MockAssignmentHistoryRepository)
	taskService := service.NewTaskService(taskRepo, columnRepo, userRepo, assignmentHistoryRepo, nil)

	taskID := uuid.New()
	oldAssigneeID := uuid.New()
	changedBy := uuid.New()

	task := &domain.Task{
		ID:         taskID,
		AssigneeID: &oldAssigneeID,
		UpdatedAt:  time.Now(),
	}

	taskRepo.On("FindByID", taskID).Return(task, nil)
	assignmentHistoryRepo.On("Create", mock.AnythingOfType("*domain.AssignmentHistory")).Return(nil)
	taskRepo.On("Update", mock.AnythingOfType("*domain.Task")).Return(nil)

	err := taskService.UpdateAssignee(taskID, nil, changedBy)

	assert.NoError(t, err)
	taskRepo.AssertExpectations(t)
	assignmentHistoryRepo.AssertExpectations(t)
}

func TestTaskService_UpdateAssignee_NoChange(t *testing.T) {
	taskRepo := new(MockTaskRepository)
	columnRepo := new(MockColumnRepository)
	userRepo := new(MockUserRepository)
	assignmentHistoryRepo := new(MockAssignmentHistoryRepository)
	taskService := service.NewTaskService(taskRepo, columnRepo, userRepo, assignmentHistoryRepo, nil)

	taskID := uuid.New()
	assigneeID := uuid.New()
	changedBy := uuid.New()

	task := &domain.Task{
		ID:         taskID,
		AssigneeID: &assigneeID,
		UpdatedAt:  time.Now(),
	}
	assignee := &domain.User{
		ID: assigneeID,
	}

	taskRepo.On("FindByID", taskID).Return(task, nil)
	userRepo.On("FindByID", assigneeID).Return(assignee, nil)
	// No history creation, no update

	err := taskService.UpdateAssignee(taskID, &assigneeID, changedBy)

	assert.NoError(t, err)
	taskRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
	assignmentHistoryRepo.AssertNotCalled(t, "Create", mock.Anything)
	taskRepo.AssertNotCalled(t, "Update", mock.Anything)
}

func TestTaskService_UpdateAssignee_TaskNotFound(t *testing.T) {
	taskRepo := new(MockTaskRepository)
	columnRepo := new(MockColumnRepository)
	userRepo := new(MockUserRepository)
	assignmentHistoryRepo := new(MockAssignmentHistoryRepository)
	taskService := service.NewTaskService(taskRepo, columnRepo, userRepo, assignmentHistoryRepo, nil)

	taskID := uuid.New()
	assigneeID := uuid.New()
	changedBy := uuid.New()

	taskRepo.On("FindByID", taskID).Return(nil, nil)

	err := taskService.UpdateAssignee(taskID, &assigneeID, changedBy)

	assert.Error(t, err)
	assert.Equal(t, service.ErrTaskNotFound, err)
	taskRepo.AssertExpectations(t)
}

func TestTaskService_UpdateAssignee_AssigneeNotFound(t *testing.T) {
	taskRepo := new(MockTaskRepository)
	columnRepo := new(MockColumnRepository)
	userRepo := new(MockUserRepository)
	assignmentHistoryRepo := new(MockAssignmentHistoryRepository)
	taskService := service.NewTaskService(taskRepo, columnRepo, userRepo, assignmentHistoryRepo, nil)

	taskID := uuid.New()
	assigneeID := uuid.New()
	changedBy := uuid.New()

	task := &domain.Task{
		ID:         taskID,
		AssigneeID: nil,
	}
	taskRepo.On("FindByID", taskID).Return(task, nil)
	userRepo.On("FindByID", assigneeID).Return(nil, nil)

	err := taskService.UpdateAssignee(taskID, &assigneeID, changedBy)

	assert.Error(t, err)
	assert.Equal(t, service.ErrAssigneeNotFound, err)
	taskRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestTaskService_GetAssignmentHistory_Success(t *testing.T) {
	taskRepo := new(MockTaskRepository)
	columnRepo := new(MockColumnRepository)
	userRepo := new(MockUserRepository)
	assignmentHistoryRepo := new(MockAssignmentHistoryRepository)
	taskService := service.NewTaskService(taskRepo, columnRepo, userRepo, assignmentHistoryRepo, nil)

	taskID := uuid.New()
	newAssigneeID := uuid.New()
	changedByID := uuid.New()
	task := &domain.Task{
		ID: taskID,
	}
	histories := []*domain.AssignmentHistory{
		{
			ID:            uuid.New(),
			TaskID:        taskID,
			OldAssigneeID: nil,
			NewAssigneeID: &newAssigneeID,
			ChangedBy:     changedByID,
			CreatedAt:     time.Now(),
		},
	}

	taskRepo.On("FindByID", taskID).Return(task, nil)
	assignmentHistoryRepo.On("FindByTaskID", taskID).Return(histories, nil)

	result, err := taskService.GetAssignmentHistory(taskID)

	assert.NoError(t, err)
	assert.Equal(t, histories, result)
	taskRepo.AssertExpectations(t)
	assignmentHistoryRepo.AssertExpectations(t)
}

func TestTaskService_GetAssignmentHistory_TaskNotFound(t *testing.T) {
	taskRepo := new(MockTaskRepository)
	columnRepo := new(MockColumnRepository)
	userRepo := new(MockUserRepository)
	assignmentHistoryRepo := new(MockAssignmentHistoryRepository)
	taskService := service.NewTaskService(taskRepo, columnRepo, userRepo, assignmentHistoryRepo, nil)

	taskID := uuid.New()
	taskRepo.On("FindByID", taskID).Return(nil, nil)

	result, err := taskService.GetAssignmentHistory(taskID)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, service.ErrTaskNotFound, err)
	taskRepo.AssertExpectations(t)
}
