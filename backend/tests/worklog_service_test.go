package tests

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"auth-project/internal/domain"
	"auth-project/internal/service"
)

// MockWorklogRepository is a mock implementation of WorklogRepository
type MockWorklogRepository struct {
	mock.Mock
}

func (m *MockWorklogRepository) Create(worklog *domain.Worklog) error {
	args := m.Called(worklog)
	return args.Error(0)
}

func (m *MockWorklogRepository) FindByTaskID(taskID uuid.UUID) ([]*domain.Worklog, error) {
	args := m.Called(taskID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Worklog), args.Error(1)
}

func (m *MockWorklogRepository) GetTimeReport() ([]*domain.TimeReportRow, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.TimeReportRow), args.Error(1)
}

func TestWorklogService_LogWork_Success(t *testing.T) {
	mockWorklogRepo := new(MockWorklogRepository)
	mockTaskRepo := new(MockTaskRepository)

	svc := service.NewWorklogService(mockWorklogRepo, mockTaskRepo)

	taskID := uuid.New()
	userID := uuid.New()
	req := service.LogWorkRequest{
		TimeSpent:   2.5,
		Description: "Implemented feature",
	}

	// Expect task exists
	task := &domain.Task{ID: taskID}
	mockTaskRepo.On("FindByID", taskID).Return(task, nil)

	// Expect worklog creation
	mockWorklogRepo.On("Create", mock.AnythingOfType("*domain.Worklog")).Return(nil)

	resp, err := svc.LogWork(taskID, userID, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, taskID, resp.TaskID)
	assert.Equal(t, userID, resp.UserID)
	assert.Equal(t, 2.5, resp.TimeSpent)
	mockTaskRepo.AssertExpectations(t)
	mockWorklogRepo.AssertExpectations(t)
}

func TestWorklogService_LogWork_InvalidTimeSpent(t *testing.T) {
	mockWorklogRepo := new(MockWorklogRepository)
	mockTaskRepo := new(MockTaskRepository)

	svc := service.NewWorklogService(mockWorklogRepo, mockTaskRepo)

	taskID := uuid.New()
	userID := uuid.New()
	req := service.LogWorkRequest{
		TimeSpent:   0,
		Description: "No time",
	}

	resp, err := svc.LogWork(taskID, userID, req)

	assert.Error(t, err)
	assert.Equal(t, service.ErrInvalidTimeSpent, err)
	assert.Nil(t, resp)
}

func TestWorklogService_LogWork_TaskNotFound(t *testing.T) {
	mockWorklogRepo := new(MockWorklogRepository)
	mockTaskRepo := new(MockTaskRepository)

	svc := service.NewWorklogService(mockWorklogRepo, mockTaskRepo)

	taskID := uuid.New()
	userID := uuid.New()
	req := service.LogWorkRequest{
		TimeSpent:   1.0,
		Description: "Task not found",
	}

	mockTaskRepo.On("FindByID", taskID).Return((*domain.Task)(nil), nil)

	resp, err := svc.LogWork(taskID, userID, req)

	assert.Error(t, err)
	assert.Equal(t, service.ErrTaskNotFound, err)
	assert.Nil(t, resp)
	mockTaskRepo.AssertExpectations(t)
}

func TestWorklogService_GetWorklogs_Success(t *testing.T) {
	mockWorklogRepo := new(MockWorklogRepository)
	mockTaskRepo := new(MockTaskRepository)

	svc := service.NewWorklogService(mockWorklogRepo, mockTaskRepo)

	taskID := uuid.New()
	userID := uuid.New()
	worklogs := []*domain.Worklog{
		{ID: uuid.New(), TaskID: taskID, UserID: userID, TimeSpent: 2.0, Description: "First"},
		{ID: uuid.New(), TaskID: taskID, UserID: userID, TimeSpent: 1.5, Description: "Second"},
	}

	task := &domain.Task{ID: taskID}
	mockTaskRepo.On("FindByID", taskID).Return(task, nil)
	mockWorklogRepo.On("FindByTaskID", taskID).Return(worklogs, nil)

	result, err := svc.GetWorklogs(taskID)

	assert.NoError(t, err)
	assert.Equal(t, worklogs, result)
	mockTaskRepo.AssertExpectations(t)
	mockWorklogRepo.AssertExpectations(t)
}

func TestWorklogService_GetWorklogs_TaskNotFound(t *testing.T) {
	mockWorklogRepo := new(MockWorklogRepository)
	mockTaskRepo := new(MockTaskRepository)

	svc := service.NewWorklogService(mockWorklogRepo, mockTaskRepo)

	taskID := uuid.New()
	mockTaskRepo.On("FindByID", taskID).Return((*domain.Task)(nil), nil)

	result, err := svc.GetWorklogs(taskID)

	assert.Error(t, err)
	assert.Equal(t, service.ErrTaskNotFound, err)
	assert.Nil(t, result)
	mockTaskRepo.AssertExpectations(t)
}

func TestWorklogService_GetTimeReport_Success(t *testing.T) {
	mockWorklogRepo := new(MockWorklogRepository)
	mockTaskRepo := new(MockTaskRepository)

	svc := service.NewWorklogService(mockWorklogRepo, mockTaskRepo)

	rows := []*domain.TimeReportRow{
		{
			TaskID:        uuid.New(),
			Title:         "Task 1",
			ColumnName:    "Todo",
			AssigneeID:    nil,
			AssigneeEmail: nil,
			TotalHours:    5.0,
		},
		{
			TaskID:        uuid.New(),
			Title:         "Task 2",
			ColumnName:    "In Progress",
			AssigneeID:    ptr(uuid.New()),
			AssigneeEmail: ptr("user@example.com"),
			TotalHours:    3.5,
		},
	}

	mockWorklogRepo.On("GetTimeReport").Return(rows, nil)

	report, err := svc.GetTimeReport()

	assert.NoError(t, err)
	assert.NotNil(t, report)
	assert.Equal(t, rows, report.Tasks)
	assert.Equal(t, 8.5, report.GrandTotal)
	mockWorklogRepo.AssertExpectations(t)
}

func TestWorklogService_GetTimeReport_Empty(t *testing.T) {
	mockWorklogRepo := new(MockWorklogRepository)
	mockTaskRepo := new(MockTaskRepository)

	svc := service.NewWorklogService(mockWorklogRepo, mockTaskRepo)

	mockWorklogRepo.On("GetTimeReport").Return([]*domain.TimeReportRow{}, nil)

	report, err := svc.GetTimeReport()

	assert.NoError(t, err)
	assert.NotNil(t, report)
	assert.Empty(t, report.Tasks)
	assert.Equal(t, 0.0, report.GrandTotal)
	mockWorklogRepo.AssertExpectations(t)
}

func TestWorklogService_GetTimeReport_Error(t *testing.T) {
	mockWorklogRepo := new(MockWorklogRepository)
	mockTaskRepo := new(MockTaskRepository)

	svc := service.NewWorklogService(mockWorklogRepo, mockTaskRepo)

	mockWorklogRepo.On("GetTimeReport").Return(([]*domain.TimeReportRow)(nil), assert.AnError)

	report, err := svc.GetTimeReport()

	assert.Error(t, err)
	assert.Nil(t, report)
	mockWorklogRepo.AssertExpectations(t)
}

// helper to create pointer
func ptr[T any](v T) *T {
	return &v
}
