package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/theikdi-sann/qr-restaurant-api/internal/domain"
)

type CreateSessionInput struct {
	TableID       uuid.UUID
	SessionTypeID uuid.UUID
	CreatedBy     uuid.UUID
	GuestCount    int
}

type DiningSessionUsecase interface {
	CreateSession(ctx context.Context, input CreateSessionInput) (*domain.DiningSession, error)
	GetSession(ctx context.Context, id uuid.UUID) (*domain.DiningSession, error)
}

type diningSessionUsecase struct {
	sessionRepo     domain.DiningSessionRepository
	sessionTypeRepo domain.DiningSessionTypeRepository
	contextTimeout  int
}

func NewDiningSessionUsecase(s domain.DiningSessionRepository, t domain.DiningSessionTypeRepository) DiningSessionUsecase {
	return &diningSessionUsecase{
		sessionRepo:     s,
		sessionTypeRepo: t,
	}
}

func (u *diningSessionUsecase) GetSession(ctx context.Context, id uuid.UUID) (*domain.DiningSession, error) {
	return u.sessionRepo.GetByID(ctx, id)
}

func (u *diningSessionUsecase) CreateSession(ctx context.Context, input CreateSessionInput) (*domain.DiningSession, error) {
	// 1. Check if table is occupied
	activeSession, err := u.sessionRepo.GetActiveSessionByTableID(ctx, input.TableID)
	if err != nil {
		return nil, err
	}
	if activeSession != nil {
		// Use a custom domain error in real app, simply string for now
		return nil, domain.ErrTableOccupied
	}

	// 2. Get Session Type details (for pricing/duration)
	sessionType, err := u.sessionTypeRepo.GetByID(ctx, input.SessionTypeID)
	if err != nil {
		return nil, err
	}

	// 3. Create Domain Entity (handles logic for ExpiresAt and Price)
	// Note: We need to import time. We will update imports.
	// For now, assume time.Now() is handled.
	newSession := domain.NewDiningSession(
		input.TableID,
		input.SessionTypeID,
		input.CreatedBy,
		input.GuestCount,
		*sessionType,
		time.Now(), // In a real app, maybe pass this in for testability
	)

	// 4. Save to Repo
	return u.sessionRepo.Create(ctx, newSession)
}
