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

// MockColumnRepository is a mock implementation of ColumnRepository
type MockColumnRepository struct {
	mock.Mock
}

func (m *MockColumnRepository) FindAll() ([]*domain.Column, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Column), args.Error(1)
}

func (m *MockColumnRepository) FindByID(id uuid.UUID) (*domain.Column, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Column), args.Error(1)
}

func (m *MockColumnRepository) FindByName(name string) (*domain.Column, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Column), args.Error(1)
}

// MockTaskRepository is a mock implementation of TaskRepository
type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(task *domain.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockTaskRepository) FindByID(id uuid.UUID) (*domain.Task, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Task), args.Error(1)
}

func (m *MockTaskRepository) FindByColumn(columnID uuid.UUID) ([]*domain.Task, error) {
	args := m.Called(columnID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Task), args.Error(1)
}

func (m *MockTaskRepository) FindAllWithColumns() ([]*domain.Task, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Task), args.Error(1)
}

func (m *MockTaskRepository) Update(task *domain.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockTaskRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTaskRepository) GetMaxPosition(columnID uuid.UUID) (float64, error) {
	args := m.Called(columnID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockTaskRepository) FindTaskBefore(columnID uuid.UUID, position float64) (*domain.Task, error) {
	args := m.Called(columnID, position)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Task), args.Error(1)
}

func (m *MockTaskRepository) FindTaskAfter(columnID uuid.UUID, position float64) (*domain.Task, error) {
	args := m.Called(columnID, position)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Task), args.Error(1)
}

func TestBoardService_GetBoard_Success(t *testing.T) {
	// Setup mocks
	columnRepo := new(MockColumnRepository)
	taskRepo := new(MockTaskRepository)
	boardService := service.NewBoardService(columnRepo, taskRepo)

	// Mock data
	column1 := &domain.Column{
		ID:        uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		Name:      "Backlog",
		Order:     0,
		CreatedAt: time.Now(),
	}
	column2 := &domain.Column{
		ID:        uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		Name:      "To Do",
		Order:     1,
		CreatedAt: time.Now(),
	}
	columns := []*domain.Column{column1, column2}

	task1 := &domain.Task{
		ID:          uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
		Title:       "Task 1",
		Description: "Description 1",
		ColumnID:    column1.ID,
		AssigneeID:  nil,
		DueDate:     nil,
		CreatedBy:   uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"),
		Position:    0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	task2 := &domain.Task{
		ID:          uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc"),
		Title:       "Task 2",
		Description: "Description 2",
		ColumnID:    column2.ID,
		AssigneeID:  nil,
		DueDate:     nil,
		CreatedBy:   uuid.MustParse("dddddddd-dddd-dddd-dddd-dddddddddddd"),
		Position:    0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	tasks := []*domain.Task{task1, task2}

	// Expectations
	columnRepo.On("FindAll").Return(columns, nil)
	taskRepo.On("FindAllWithColumns").Return(tasks, nil)

	// Execute
	board, err := boardService.GetBoard()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, board)
	assert.Len(t, board.Columns, 2)
	assert.Equal(t, "Backlog", board.Columns[0].Name)
	assert.Equal(t, "To Do", board.Columns[1].Name)
	assert.Len(t, board.Columns[0].Tasks, 1)
	assert.Len(t, board.Columns[1].Tasks, 1)
	assert.Equal(t, "Task 1", board.Columns[0].Tasks[0].Title)
	assert.Equal(t, "Task 2", board.Columns[1].Tasks[0].Title)

	columnRepo.AssertExpectations(t)
	taskRepo.AssertExpectations(t)
}

func TestBoardService_GetBoard_ColumnRepoError(t *testing.T) {
	columnRepo := new(MockColumnRepository)
	taskRepo := new(MockTaskRepository)
	boardService := service.NewBoardService(columnRepo, taskRepo)

	columnRepo.On("FindAll").Return(nil, assert.AnError)

	board, err := boardService.GetBoard()

	assert.Error(t, err)
	assert.Nil(t, board)
	columnRepo.AssertExpectations(t)
	taskRepo.AssertExpectations(t)
}

func TestBoardService_GetBoard_TaskRepoError(t *testing.T) {
	columnRepo := new(MockColumnRepository)
	taskRepo := new(MockTaskRepository)
	boardService := service.NewBoardService(columnRepo, taskRepo)

	columns := []*domain.Column{
		{ID: uuid.New(), Name: "Backlog", Order: 0, CreatedAt: time.Now()},
	}
	columnRepo.On("FindAll").Return(columns, nil)
	taskRepo.On("FindAllWithColumns").Return(nil, assert.AnError)

	board, err := boardService.GetBoard()

	assert.Error(t, err)
	assert.Nil(t, board)
	columnRepo.AssertExpectations(t)
	taskRepo.AssertExpectations(t)
}

func TestBoardService_GetBoard_EmptyBoard(t *testing.T) {
	columnRepo := new(MockColumnRepository)
	taskRepo := new(MockTaskRepository)
	boardService := service.NewBoardService(columnRepo, taskRepo)

	columns := []*domain.Column{
		{ID: uuid.New(), Name: "Backlog", Order: 0, CreatedAt: time.Now()},
		{ID: uuid.New(), Name: "To Do", Order: 1, CreatedAt: time.Now()},
	}
	tasks := []*domain.Task{}

	columnRepo.On("FindAll").Return(columns, nil)
	taskRepo.On("FindAllWithColumns").Return(tasks, nil)

	board, err := boardService.GetBoard()

	assert.NoError(t, err)
	assert.NotNil(t, board)
	assert.Len(t, board.Columns, 2)
	assert.Empty(t, board.Columns[0].Tasks)
	assert.Empty(t, board.Columns[1].Tasks)

	columnRepo.AssertExpectations(t)
	taskRepo.AssertExpectations(t)
}
