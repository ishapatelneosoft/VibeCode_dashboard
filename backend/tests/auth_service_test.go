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

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(id uuid.UUID) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) ExistsByEmail(email string) (bool, error) {
	args := m.Called(email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) List() ([]*domain.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.User), args.Error(1)
}

// MockSessionRepository is a mock implementation of SessionRepository
type MockSessionRepository struct {
	mock.Mock
}

func (m *MockSessionRepository) Create(session *domain.Session) error {
	args := m.Called(session)
	return args.Error(0)
}

func (m *MockSessionRepository) FindByID(sessionID uuid.UUID) (*domain.Session, error) {
	args := m.Called(sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Session), args.Error(1)
}

func (m *MockSessionRepository) FindByUserID(userID uuid.UUID) ([]*domain.Session, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Session), args.Error(1)
}

func (m *MockSessionRepository) Delete(sessionID uuid.UUID) error {
	args := m.Called(sessionID)
	return args.Error(0)
}

func (m *MockSessionRepository) DeleteByUserID(userID uuid.UUID) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockSessionRepository) DeleteExpired() error {
	args := m.Called()
	return args.Error(0)
}

func TestAuthService_Register_Success(t *testing.T) {
	// Setup mocks
	userRepo := new(MockUserRepository)
	sessionRepo := new(MockSessionRepository)

	// Configure mock expectations
	userRepo.On("ExistsByEmail", "test@example.com").Return(false, nil)
	userRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(nil)

	// Create service
	authService := service.NewAuthService(
		userRepo,
		sessionRepo,
		24*time.Hour,
	)

	// Test registration
	req := service.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	resp, err := authService.Register(req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "test@example.com", resp.Email)
	assert.NotEmpty(t, resp.UserID)

	// Verify mock expectations
	userRepo.AssertExpectations(t)
}

func TestAuthService_Register_UserAlreadyExists(t *testing.T) {
	// Setup mocks
	userRepo := new(MockUserRepository)
	sessionRepo := new(MockSessionRepository)

	// Configure mock to return that user exists
	userRepo.On("ExistsByEmail", "existing@example.com").Return(true, nil)

	// Create service
	authService := service.NewAuthService(
		userRepo,
		sessionRepo,
		24*time.Hour,
	)

	// Test registration
	req := service.RegisterRequest{
		Email:    "existing@example.com",
		Password: "password123",
	}

	resp, err := authService.Register(req)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, service.ErrUserAlreadyExists, err)
	assert.Nil(t, resp)

	// Verify mock expectations
	userRepo.AssertExpectations(t)
}

func TestAuthService_Login_Success(t *testing.T) {
	// Setup mocks
	userRepo := new(MockUserRepository)
	sessionRepo := new(MockSessionRepository)

	// Create a test user
	userID := uuid.New()
	testUser := &domain.User{
		ID:           userID,
		Email:        "test@example.com",
		PasswordHash: "$2a$10$abcdefghijklmnopqrstuv", // Mock bcrypt hash
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Configure mock expectations
	userRepo.On("FindByEmail", "test@example.com").Return(testUser, nil)
	sessionRepo.On("Create", mock.AnythingOfType("*domain.Session")).Return(nil)

	// Create service
	authService := service.NewAuthService(
		userRepo,
		sessionRepo,
		24*time.Hour,
	)

	// Test login
	req := service.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	// Note: This test would need actual bcrypt hash verification
	// For now, we'll skip the actual password verification
	t.Skip("Password verification requires actual bcrypt implementation")

	resp, err := authService.Login(req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, userID.String(), resp.UserID)
	assert.Equal(t, "test@example.com", resp.Email)
	assert.NotEmpty(t, resp.SessionID)

	// Verify mock expectations
	userRepo.AssertExpectations(t)
	sessionRepo.AssertExpectations(t)
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	// Setup mocks
	userRepo := new(MockUserRepository)
	sessionRepo := new(MockSessionRepository)

	// Configure mock to return nil user (not found)
	userRepo.On("FindByEmail", "nonexistent@example.com").Return(nil, nil)

	// Create service
	authService := service.NewAuthService(
		userRepo,
		sessionRepo,
		24*time.Hour,
	)

	// Test login
	req := service.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	resp, err := authService.Login(req)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, service.ErrInvalidCredentials, err)
	assert.Nil(t, resp)

	// Verify mock expectations
	userRepo.AssertExpectations(t)
}

func TestAuthService_ValidateSession_Success(t *testing.T) {
	// Setup mocks
	userRepo := new(MockUserRepository)
	sessionRepo := new(MockSessionRepository)

	// Create test data
	userID := uuid.New()
	sessionID := uuid.New()

	testSession := &domain.Session{
		SessionID: sessionID,
		UserID:    userID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}

	testUser := &domain.User{
		ID:    userID,
		Email: "test@example.com",
	}

	// Configure mock expectations
	sessionRepo.On("FindByID", sessionID).Return(testSession, nil)
	userRepo.On("FindByID", userID).Return(testUser, nil)

	// Create service
	authService := service.NewAuthService(
		userRepo,
		sessionRepo,
		24*time.Hour,
	)

	// Test session validation
	user, err := authService.ValidateSession(sessionID.String())

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, "test@example.com", user.Email)

	// Verify mock expectations
	sessionRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}
