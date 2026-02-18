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
	ListActiveSessions(ctx context.Context) ([]*domain.DiningSession, error)
	ExtendSession(ctx context.Context, id uuid.UUID, minutes int) (*domain.DiningSession, error)
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

func (u *diningSessionUsecase) ListActiveSessions(ctx context.Context) ([]*domain.DiningSession, error) {
	return u.sessionRepo.ListActiveSessions(ctx)
}

func (u *diningSessionUsecase) GetSession(ctx context.Context, id uuid.UUID) (*domain.DiningSession, error) {
	return u.sessionRepo.GetByID(ctx, id)
}

func (u *diningSessionUsecase) ExtendSession(ctx context.Context, id uuid.UUID, minutes int) (*domain.DiningSession, error) {
	session, err := u.sessionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, domain.ErrNotFound
	}

	var newExpiry time.Time
	if session.ExpiresAt != nil && session.ExpiresAt.After(time.Now()) {
		newExpiry = session.ExpiresAt.Add(time.Duration(minutes) * time.Minute)
	} else {
		newExpiry = time.Now().Add(time.Duration(minutes) * time.Minute)
	}

	return u.sessionRepo.Extend(ctx, id, newExpiry)
}

func (u *diningSessionUsecase) CreateSession(ctx context.Context, input CreateSessionInput) (*domain.DiningSession, error) {
	// 1. Check if table is occupied (if table is assigned)
	if input.TableID != uuid.Nil {
		activeSession, err := u.sessionRepo.GetActiveSessionByTableID(ctx, input.TableID)
		if err != nil {
			return nil, err
		}
		if activeSession != nil {
			return nil, domain.ErrTableOccupied
		}
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
