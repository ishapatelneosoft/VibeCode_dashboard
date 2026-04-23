package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"auth-project/internal/domain"
)

// SessionRepositoryPostgres implements SessionRepository for PostgreSQL
type SessionRepositoryPostgres struct {
	db *pgxpool.Pool
}

// NewSessionRepositoryPostgres creates a new PostgreSQL session repository
func NewSessionRepositoryPostgres(db *pgxpool.Pool) *SessionRepositoryPostgres {
	return &SessionRepositoryPostgres{db: db}
}

// Create inserts a new session into the database
func (r *SessionRepositoryPostgres) Create(session *domain.Session) error {
	query := `
		INSERT INTO sessions (session_id, user_id, expires_at, created_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(
		context.Background(),
		query,
		session.SessionID,
		session.UserID,
		session.ExpiresAt,
		session.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

// FindByID retrieves a session by its ID
func (r *SessionRepositoryPostgres) FindByID(sessionID uuid.UUID) (*domain.Session, error) {
	query := `
		SELECT session_id, user_id, expires_at, created_at
		FROM sessions
		WHERE session_id = $1
	`

	var session domain.Session
	err := r.db.QueryRow(
		context.Background(),
		query,
		sessionID,
	).Scan(
		&session.SessionID,
		&session.UserID,
		&session.ExpiresAt,
		&session.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find session by ID: %w", err)
	}

	return &session, nil
}

// FindByUserID retrieves all sessions for a user
func (r *SessionRepositoryPostgres) FindByUserID(userID uuid.UUID) ([]*domain.Session, error) {
	query := `
		SELECT session_id, user_id, expires_at, created_at
		FROM sessions
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(context.Background(), query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query sessions: %w", err)
	}
	defer rows.Close()

	var sessions []*domain.Session
	for rows.Next() {
		var session domain.Session
		if err := rows.Scan(
			&session.SessionID,
			&session.UserID,
			&session.ExpiresAt,
			&session.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}
		sessions = append(sessions, &session)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating sessions: %w", err)
	}

	return sessions, nil
}

// Delete removes a session by ID
func (r *SessionRepositoryPostgres) Delete(sessionID uuid.UUID) error {
	query := `DELETE FROM sessions WHERE session_id = $1`

	result, err := r.db.Exec(context.Background(), query, sessionID)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("session not found")
	}

	return nil
}

// DeleteByUserID removes all sessions for a user
func (r *SessionRepositoryPostgres) DeleteByUserID(userID uuid.UUID) error {
	query := `DELETE FROM sessions WHERE user_id = $1`

	_, err := r.db.Exec(context.Background(), query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user sessions: %w", err)
	}

	return nil
}

// DeleteExpired removes all expired sessions
func (r *SessionRepositoryPostgres) DeleteExpired() error {
	query := `DELETE FROM sessions WHERE expires_at < $1`

	_, err := r.db.Exec(context.Background(), query, time.Now())
	if err != nil {
		return fmt.Errorf("failed to delete expired sessions: %w", err)
	}

	return nil
}
