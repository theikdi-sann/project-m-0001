package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/theikdi-sann/qr-restaurant-api/internal/domain"
	"github.com/theikdi-sann/qr-restaurant-api/internal/repository/postgres/db"
)

type diningSessionRepository struct {
	pool *pgxpool.Pool
}

func NewDiningSessionRepository(pool *pgxpool.Pool) domain.DiningSessionRepository {
	return &diningSessionRepository{
		pool: pool,
	}
}

func (r *diningSessionRepository) GetActiveSessionByTableID(ctx context.Context, tableID uuid.UUID) (*domain.DiningSession, error) {
	row, err := db.New(r.pool).GetActiveSessionByTableID(ctx, uuidToPg(tableID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // No active session found
		}
		return nil, err
	}
	return mapToDomain(row), nil
}

func (r *diningSessionRepository) Create(ctx context.Context, session *domain.DiningSession) (*domain.DiningSession, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := db.New(r.pool).WithTx(tx)

	// 1. Create Session
	arg := db.CreateDiningSessionParams{
		TableID:       uuidToPg(session.TableID),
		SessionTypeID: uuidToPg(session.SessionTypeID),
		CreatedBy:     uuidToPg(session.CreatedBy),
		GuestCount:    int32(session.GuestCount),
		Status:        string(session.Status),
		StartTime:     timeToPg(session.StartTime),
		ExpiresAt:     timePtrToPg(session.ExpiresAt),
		PricePerGuest: decimalToPg(session.PricePerGuest),
		TotalAmount:   decimalToPg(session.TotalAmount),
	}

	row, err := qtx.CreateDiningSession(ctx, arg)
	if err != nil {
		return nil, err
	}

	// 2. Update Table Status to Occupied
	err = qtx.UpdateTableStatus(ctx, db.UpdateTableStatusParams{
		ID:     arg.TableID,
		Status: "occupied",
	})
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return mapToDomain(row), nil
}

func (r *diningSessionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.DiningSession, error) {
	row, err := db.New(r.pool).GetDiningSession(ctx, uuidToPg(id))
	if err != nil {
		return nil, err
	}
	return mapToDomain(row), nil
}

func (r *diningSessionRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.SessionStatus) (*domain.DiningSession, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	qtx := db.New(r.pool).WithTx(tx)

	// 1. Update Session Status
	arg := db.UpdateDiningSessionStatusParams{
		ID:     uuidToPg(id),
		Status: string(status),
	}
	row, err := qtx.UpdateDiningSessionStatus(ctx, arg)
	if err != nil {
		return nil, err
	}

	// 2. Update Table Status if needed
	// Completed or Cancelled -> Available
	// Expired -> Occupied (Guests still there)
	if status == domain.SessionStatusCompleted || status == domain.SessionStatusCancelled {
		err = qtx.UpdateTableStatus(ctx, db.UpdateTableStatusParams{
			ID:     row.TableID,
			Status: "available",
		})
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return mapToDomain(row), nil
}

func (r *diningSessionRepository) ListExpiredActiveSessions(ctx context.Context) ([]*domain.DiningSession, error) {
	rows, err := db.New(r.pool).ListExpiredSessions(ctx)
	if err != nil {
		return nil, err
	}

	var sessions []*domain.DiningSession
	for _, row := range rows {
		sessions = append(sessions, mapToDomain(row))
	}
	return sessions, nil
}

func (r *diningSessionRepository) Extend(ctx context.Context, id uuid.UUID, newExpiry time.Time) (*domain.DiningSession, error) {
	arg := db.ExtendSessionParams{
		ID:        uuidToPg(id),
		ExpiresAt: timeToPg(newExpiry),
	}
	row, err := db.New(r.pool).ExtendSession(ctx, arg)
	if err != nil {
		return nil, err
	}
	return mapToDomain(row), nil
}

func mapToDomain(row db.DiningSession) *domain.DiningSession {
	return &domain.DiningSession{
		ID:            pgToUuid(row.ID),
		TableID:       pgToUuid(row.TableID),
		SessionTypeID: pgToUuid(row.SessionTypeID),
		CreatedBy:     pgToUuid(row.CreatedBy),
		GuestCount:    int(row.GuestCount),
		Status:        domain.SessionStatus(row.Status),
		StartTime:     row.StartTime.Time,
		ExpiresAt:     pgTimePtr(row.ExpiresAt),
		PricePerGuest: pgToDecimal(row.PricePerGuest),
		TotalAmount:   pgToDecimal(row.TotalAmount),
	}
}
