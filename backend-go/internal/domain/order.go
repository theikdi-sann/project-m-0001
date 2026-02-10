package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPreparing OrderStatus = "preparing"
	OrderStatusServed    OrderStatus = "served"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	ID              uuid.UUID
	DiningSessionID uuid.UUID
	Items           []OrderItem
	Status          OrderStatus
	TotalAmount     decimal.Decimal
	CreatedAt       time.Time
}

type OrderItem struct {
	ID         uuid.UUID
	OrderID    uuid.UUID
	MenuItemID uuid.UUID
	Quantity   int
	UnitPrice  decimal.Decimal
	Notes      string
}

// ValidateOrder checks if the order is valid based on the session status.
func ValidateOrder(session *DiningSession, order *Order) error {
	if session.Status != SessionStatusActive {
		return ErrSessionNotActive
	}

	if session.ExpiresAt != nil && time.Now().After(*session.ExpiresAt) {
		return ErrSessionExpired
	}

	return nil
}