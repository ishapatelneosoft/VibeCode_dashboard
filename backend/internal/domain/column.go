package domain

import (
	"time"

	"github.com/google/uuid"
)

// Column represents a board column entity
type Column struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Order     int       `json:"order"`
	CreatedAt time.Time `json:"created_at"`
}

// ColumnRepository defines the interface for column data operations
type ColumnRepository interface {
	FindAll() ([]*Column, error)
	FindByID(id uuid.UUID) (*Column, error)
	FindByName(name string) (*Column, error)
}
