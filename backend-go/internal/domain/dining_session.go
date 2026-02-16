package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type SessionStatus string

const (
	SessionStatusActive    SessionStatus = "active"
	SessionStatusCompleted SessionStatus = "completed"
	SessionStatusCancelled SessionStatus = "cancelled"
	SessionStatusExpired   SessionStatus = "expired"
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
	PricePerGuest decimal.Decimal
	TotalAmount   decimal.Decimal
}

type DiningSessionRepository interface {
	Create(ctx context.Context, session *DiningSession) (*DiningSession, error)
	GetByID(ctx context.Context, id uuid.UUID) (*DiningSession, error)
	GetActiveSessionByTableID(ctx context.Context, tableID uuid.UUID) (*DiningSession, error)
	ListExpiredActiveSessions(ctx context.Context) ([]*DiningSession, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status SessionStatus) (*DiningSession, error)
	Extend(ctx context.Context, id uuid.UUID, newExpiry time.Time) (*DiningSession, error)
}

type DiningSessionTypeRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*DiningSessionType, error)
}

type DiningSessionType struct {
	ID              uuid.UUID
	Name            string
	IsBuffet        bool
	Price           decimal.Decimal
	DurationMinutes int
}

// NewDiningSession creates a new dining session and calculates the expiry time if it's a buffet.
func NewDiningSession(tableID, sessionTypeID, createdBy uuid.UUID, guestCount int, sessionType DiningSessionType, startTime time.Time) *DiningSession {
	var expiresAt *time.Time
	if sessionType.IsBuffet && sessionType.DurationMinutes > 0 {
		expiryTime := startTime.Add(time.Duration(sessionType.DurationMinutes) * time.Minute)
		expiresAt = &expiryTime
	}

	pricePerGuest := decimal.Zero
	totalAmount := decimal.Zero
	if sessionType.IsBuffet {
		pricePerGuest = sessionType.Price
		totalAmount = pricePerGuest.Mul(decimal.NewFromInt(int64(guestCount)))
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