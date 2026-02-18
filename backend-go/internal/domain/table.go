package domain

import (
	"context"

	"github.com/google/uuid"
)

type DiningTable struct {
	ID          uuid.UUID
	TableNumber string
	Capacity    int
	Status      string // available, occupied
}

type TableRepository interface {
	ListAll(ctx context.Context) ([]*DiningTable, error)
	GetByID(ctx context.Context, id uuid.UUID) (*DiningTable, error)
}
