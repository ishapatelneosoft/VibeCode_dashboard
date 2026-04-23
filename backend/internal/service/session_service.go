package service

import (
	"time"

	"github.com/google/uuid"

	"auth-project/internal/domain"
)

// SessionService handles session management business logic
type SessionService struct {
	sessionRepo domain.SessionRepository
	defaultTTL  time.Duration
}

// NewSessionService creates a new session service
func NewSessionService(sessionRepo domain.SessionRepository, defaultTTL time.Duration) *SessionService {
	return &SessionService{
		sessionRepo: sessionRepo,
		defaultTTL:  defaultTTL,
	}
}

// CreateSession creates a new session for a user
func (s *SessionService) CreateSession(userID uuid.UUID, customTTL ...time.Duration) (*domain.Session, error) {
	ttl := s.defaultTTL
	if len(customTTL) > 0 {
		ttl = customTTL[0]
	}

	session := domain.NewSession(userID, ttl)
	if err := s.sessionRepo.Create(session); err != nil {
		return nil, err
	}

	return session, nil
}

// GetSession retrieves a session by ID
func (s *SessionService) GetSession(sessionID uuid.UUID) (*domain.Session, error) {
	session, err := s.sessionRepo.FindByID(sessionID)
	if err != nil {
		return nil, err
	}

	if session == nil {
		return nil, nil
	}

	// Check if session is expired
	if session.IsExpired() {
		// Clean up expired session
		_ = s.sessionRepo.Delete(sessionID)
		return nil, nil
	}

	return session, nil
}

// GetUserSessions retrieves all active sessions for a user
func (s *SessionService) GetUserSessions(userID uuid.UUID) ([]*domain.Session, error) {
	sessions, err := s.sessionRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Filter out expired sessions
	var activeSessions []*domain.Session
	for _, session := range sessions {
		if !session.IsExpired() {
			activeSessions = append(activeSessions, session)
		}
	}

	return activeSessions, nil
}

// RevokeSession revokes a specific session
func (s *SessionService) RevokeSession(sessionID uuid.UUID) error {
	return s.sessionRepo.Delete(sessionID)
}

// RevokeAllUserSessions revokes all sessions for a user
func (s *SessionService) RevokeAllUserSessions(userID uuid.UUID) error {
	return s.sessionRepo.DeleteByUserID(userID)
}

// CleanupExpiredSessions removes all expired sessions
func (s *SessionService) CleanupExpiredSessions() error {
	return s.sessionRepo.DeleteExpired()
}

// IsSessionValid checks if a session is valid and not expired
func (s *SessionService) IsSessionValid(sessionID uuid.UUID) (bool, error) {
	session, err := s.GetSession(sessionID)
	if err != nil {
		return false, err
	}

	return session != nil, nil
}
