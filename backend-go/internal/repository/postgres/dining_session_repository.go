package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/theikdi-sann/qr-restaurant-api/internal/domain"
	"github.com/theikdi-sann/qr-restaurant-api/internal/repository/postgres/db"
)

type diningSessionRepository struct {
	queries *db.Queries
}

func NewDiningSessionRepository(queries *db.Queries) domain.DiningSessionRepository {
	return &diningSessionRepository{
		queries: queries,
	}
}

func (r *diningSessionRepository) GetActiveSessionByTableID(ctx context.Context, tableID uuid.UUID) (*domain.DiningSession, error) {
	row, err := r.queries.GetActiveSessionByTableID(ctx, tableID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // No active session found
		}
		return nil, err
	}
	return mapToDomain(row), nil
}

func (r *diningSessionRepository) Create(ctx context.Context, session *domain.DiningSession) (*domain.DiningSession, error) {
	arg := db.CreateDiningSessionParams{
		TableID:       session.TableID,
		SessionTypeID: session.SessionTypeID,
		CreatedBy:     session.CreatedBy,
		GuestCount:    int32(session.GuestCount),
		Status:        string(session.Status),
		StartTime:     session.StartTime,
		ExpiresAt:     session.ExpiresAt,
		PricePerGuest: session.PricePerGuest,
		TotalAmount:   session.TotalAmount,
	}

	row, err := r.queries.CreateDiningSession(ctx, arg)
	if err != nil {
		return nil, err
	}

	return mapToDomain(row), nil
}

func (r *diningSessionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.DiningSession, error) {
	row, err := r.queries.GetDiningSession(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapToDomain(row), nil
}

func (r *diningSessionRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.SessionStatus) (*domain.DiningSession, error) {
	arg := db.UpdateDiningSessionStatusParams{
		ID:     id,
		Status: string(status),
	}
	row, err := r.queries.UpdateDiningSessionStatus(ctx, arg)
	if err != nil {
		return nil, err
	}
	return mapToDomain(row), nil
}

func mapToDomain(row db.DiningSession) *domain.DiningSession {
	return &domain.DiningSession{
		ID:            row.ID,
		TableID:       row.TableID,
		SessionTypeID: row.SessionTypeID,
		CreatedBy:     row.CreatedBy,
		GuestCount:    int(row.GuestCount),
		Status:        domain.SessionStatus(row.Status),
		StartTime:     row.StartTime,
		ExpiresAt:     row.ExpiresAt,
		PricePerGuest: row.PricePerGuest,
		TotalAmount:   row.TotalAmount,
	}
}
