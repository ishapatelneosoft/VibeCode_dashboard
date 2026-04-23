package service

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"auth-project/internal/domain"
)

var (
	// ErrInvalidCredentials is returned when credentials are invalid
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrUserAlreadyExists is returned when trying to register an existing user
	ErrUserAlreadyExists = errors.New("user already exists")
	// ErrInvalidEmail is returned when email format is invalid
	ErrInvalidEmail = errors.New("invalid email format")
	// ErrWeakPassword is returned when password doesn't meet requirements
	ErrWeakPassword = errors.New("password must be at least 8 characters long")
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo    domain.UserRepository
	sessionRepo domain.SessionRepository
	sessionTTL  time.Duration
}

// NewAuthService creates a new authentication service
func NewAuthService(
	userRepo domain.UserRepository,
	sessionRepo domain.SessionRepository,
	sessionTTL time.Duration,
) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		sessionTTL:  sessionTTL,
	}
}

// RegisterRequest represents user registration request
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterResponse represents registration response
type RegisterResponse struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// Register creates a new user account
func (s *AuthService) Register(req RegisterRequest) (*RegisterResponse, error) {
	// Validate email
	if !isValidEmail(req.Email) {
		return nil, ErrInvalidEmail
	}

	// Validate password
	if len(req.Password) < 8 {
		return nil, ErrWeakPassword
	}

	// Check if user already exists
	exists, err := s.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
		return nil, ErrUserAlreadyExists
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := domain.NewUser(req.Email, string(passwordHash))
	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &RegisterResponse{
		UserID:    user.ID.String(),
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

// LoginRequest represents user login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents login response with session
type LoginResponse struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	SessionID string    `json:"session_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Login authenticates a user and creates a session
func (s *AuthService) Login(req LoginRequest) (*LoginResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Create session
	session := domain.NewSession(user.ID, s.sessionTTL)
	if err := s.sessionRepo.Create(session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &LoginResponse{
		UserID:    user.ID.String(),
		Email:     user.Email,
		SessionID: session.SessionID.String(),
		ExpiresAt: session.ExpiresAt,
	}, nil
}

// Logout invalidates a session
func (s *AuthService) Logout(sessionID string) error {
	sid, err := parseUUID(sessionID)
	if err != nil {
		return fmt.Errorf("invalid session ID: %w", err)
	}

	if err := s.sessionRepo.Delete(sid); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

// ValidateSession checks if a session is valid
func (s *AuthService) ValidateSession(sessionID string) (*domain.User, error) {
	sid, err := parseUUID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("invalid session ID: %w", err)
	}

	session, err := s.sessionRepo.FindByID(sid)
	if err != nil {
		return nil, fmt.Errorf("failed to find session: %w", err)
	}
	if session == nil {
		return nil, errors.New("session not found")
	}

	if session.IsExpired() {
		// Clean up expired session
		_ = s.sessionRepo.Delete(sid)
		return nil, errors.New("session expired")
	}

	user, err := s.userRepo.FindByID(session.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// ListUsers retrieves all users (for assignment dropdown)
func (s *AuthService) ListUsers() ([]*domain.User, error) {
	users, err := s.userRepo.List()
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	return users, nil
}

// Helper functions
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func parseUUID(str string) (uuid.UUID, error) {
	return uuid.Parse(str)
}
