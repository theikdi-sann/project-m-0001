package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type SessionStatus string

const (
	SessionStatusActive    SessionStatus = "active"
	SessionStatusCompleted SessionStatus = "completed"
	SessionStatusCancelled SessionStatus = "cancelled"
)

type DiningSession struct {
	ID            uuid.UUID
	TableID       uuid.UUID
	SessionTypeID uuid.UUID
	CreatedBy     uuid.UUID
	GuestCount    int
	Status        SessionStatus
	StartTime     time.Time
	ExpiresAt     *time.Time
	PricePerGuest float64
	TotalAmount   float64
}

type DiningSessionRepository interface {
	Create(ctx context.Context, session *DiningSession) (*DiningSession, error)
	GetByID(ctx context.Context, id uuid.UUID) (*DiningSession, error)
	GetActiveSessionByTableID(ctx context.Context, tableID uuid.UUID) (*DiningSession, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status SessionStatus) (*DiningSession, error)
}

type DiningSessionTypeRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*DiningSessionType, error)
}

type DiningSessionType struct {
	ID              uuid.UUID
	Name            string
	IsBuffet        bool
	Price           float64
	DurationMinutes int
}

// NewDiningSession creates a new dining session and calculates the expiry time if it's a buffet.
func NewDiningSession(tableID, sessionTypeID, createdBy uuid.UUID, guestCount int, sessionType DiningSessionType, startTime time.Time) *DiningSession {
	var expiresAt *time.Time
	if sessionType.IsBuffet && sessionType.DurationMinutes > 0 {
		expiryTime := startTime.Add(time.Duration(sessionType.DurationMinutes) * time.Minute)
		expiresAt = &expiryTime
	}

	pricePerGuest := 0.0
	totalAmount := 0.0
	if sessionType.IsBuffet {
		pricePerGuest = sessionType.Price
		totalAmount = pricePerGuest * float64(guestCount)
	}

	return &DiningSession{
		ID:            uuid.New(),
		TableID:       tableID,
		SessionTypeID: sessionTypeID,
		CreatedBy:     createdBy,
		GuestCount:    guestCount,
		Status:        SessionStatusActive,
		StartTime:     startTime,
		ExpiresAt:     expiresAt,
		PricePerGuest: pricePerGuest,
		TotalAmount:   totalAmount,
	}
}
