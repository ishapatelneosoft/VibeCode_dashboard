package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"auth-project/internal/domain"
)

// ColumnRepositoryPostgres implements ColumnRepository for PostgreSQL
type ColumnRepositoryPostgres struct {
	db *pgxpool.Pool
}

// NewColumnRepositoryPostgres creates a new PostgreSQL column repository
func NewColumnRepositoryPostgres(db *pgxpool.Pool) *ColumnRepositoryPostgres {
	return &ColumnRepositoryPostgres{db: db}
}

// FindAll retrieves all columns ordered by their order field
func (r *ColumnRepositoryPostgres) FindAll() ([]*domain.Column, error) {
	query := `
		SELECT id, name, "order", created_at
		FROM columns
		ORDER BY "order" ASC
	`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to query columns: %w", err)
	}
	defer rows.Close()

	var columns []*domain.Column
	for rows.Next() {
		var column domain.Column
		if err := rows.Scan(
			&column.ID,
			&column.Name,
			&column.Order,
			&column.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan column: %w", err)
		}
		columns = append(columns, &column)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating columns: %w", err)
	}

	return columns, nil
}

// FindByID retrieves a column by its ID
func (r *ColumnRepositoryPostgres) FindByID(id uuid.UUID) (*domain.Column, error) {
	query := `
		SELECT id, name, "order", created_at
		FROM columns
		WHERE id = $1
	`

	var column domain.Column
	err := r.db.QueryRow(
		context.Background(),
		query,
		id,
	).Scan(
		&column.ID,
		&column.Name,
		&column.Order,
		&column.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find column by ID: %w", err)
	}

	return &column, nil
}

// FindByName retrieves a column by its name
func (r *ColumnRepositoryPostgres) FindByName(name string) (*domain.Column, error) {
	query := `
		SELECT id, name, "order", created_at
		FROM columns
		WHERE name = $1
	`

	var column domain.Column
	err := r.db.QueryRow(
		context.Background(),
		query,
		name,
	).Scan(
		&column.ID,
		&column.Name,
		&column.Order,
		&column.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find column by name: %w", err)
	}

	return &column, nil
}
