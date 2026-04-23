package domain

import (
	"time"

	"github.com/google/uuid"
)

// Session represents a user session entity
type Session struct {
	SessionID uuid.UUID `json:"session_id"`
	UserID    uuid.UUID `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// NewSession creates a new Session instance with generated ID and timestamps
func NewSession(userID uuid.UUID, duration time.Duration) *Session {
	now := time.Now()
	return &Session{
		SessionID: uuid.New(),
		UserID:    userID,
		ExpiresAt: now.Add(duration),
		CreatedAt: now,
	}
}

// IsExpired checks if the session has expired
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// SessionRepository defines the interface for session data operations
type SessionRepository interface {
	Create(session *Session) error
	FindByID(sessionID uuid.UUID) (*Session, error)
	FindByUserID(userID uuid.UUID) ([]*Session, error)
	Delete(sessionID uuid.UUID) error
	DeleteByUserID(userID uuid.UUID) error
	DeleteExpired() error
}
