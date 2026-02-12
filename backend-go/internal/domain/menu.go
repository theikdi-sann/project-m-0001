package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type MenuItem struct {
	ID          uuid.UUID
	CategoryID  uuid.UUID
	Name        string
	Description *string
	Price       decimal.Decimal
	IsAvailable bool
}

type MenuItemRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*MenuItem, error)
	ListByCategory(ctx context.Context, categoryID uuid.UUID) ([]*MenuItem, error)
}
