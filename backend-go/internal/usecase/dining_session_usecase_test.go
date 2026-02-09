package usecase

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/theikdi-sann/qr-restaurant-api/internal/domain"
)

// MockSessionRepo
type MockSessionRepo struct {
	mock.Mock
}

func (m *MockSessionRepo) Create(ctx context.Context, session *domain.DiningSession) (*domain.DiningSession, error) {
	args := m.Called(ctx, session)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.DiningSession), args.Error(1)
}

func (m *MockSessionRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.DiningSession, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.DiningSession), args.Error(1)
}

func (m *MockSessionRepo) GetActiveSessionByTableID(ctx context.Context, tableID uuid.UUID) (*domain.DiningSession, error) {
	args := m.Called(ctx, tableID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.DiningSession), args.Error(1)
}

func (m *MockSessionRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.SessionStatus) (*domain.DiningSession, error) {
	args := m.Called(ctx, id, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.DiningSession), args.Error(1)
}

// MockSessionTypeRepo
type MockSessionTypeRepo struct {
	mock.Mock
}

func (m *MockSessionTypeRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.DiningSessionType, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.DiningSessionType), args.Error(1)
}

func TestCreateSession(t *testing.T) {
	mockSessionRepo := new(MockSessionRepo)
	mockTypeRepo := new(MockSessionTypeRepo)
	usecase := NewDiningSessionUsecase(mockSessionRepo, mockTypeRepo)

	tableID := uuid.New()
	sessionTypeID := uuid.New()
	createdBy := uuid.New()

	input := CreateSessionInput{
		TableID:       tableID,
		SessionTypeID: sessionTypeID,
		CreatedBy:     createdBy,
		GuestCount:    2,
	}

	t.Run("should fail if table is already occupied", func(t *testing.T) {
		// Mock: Table has an active session
		existingSession := &domain.DiningSession{ID: uuid.New(), Status: domain.SessionStatusActive}
		mockSessionRepo.On("GetActiveSessionByTableID", mock.Anything, tableID).Return(existingSession, nil).Once()

		_, err := usecase.CreateSession(context.Background(), input)

		assert.Error(t, err)
		assert.Equal(t, "table is already occupied", err.Error())
		mockSessionRepo.AssertExpectations(t)
	})

	t.Run("should create session successfully if table is free", func(t *testing.T) {
		// Mock: Table is free
		mockSessionRepo.On("GetActiveSessionByTableID", mock.Anything, tableID).Return(nil, nil).Once()

		// Mock: Get Session Type
		sessionType := &domain.DiningSessionType{
			ID:              sessionTypeID,
			IsBuffet:        true,
			Price:           10.0,
			DurationMinutes: 90,
		}
		mockTypeRepo.On("GetByID", mock.Anything, sessionTypeID).Return(sessionType, nil).Once()

		// Mock: Create Session
		// We use mock.MatchedBy to validate the session passed to repo has correct calculated fields
		mockSessionRepo.On("Create", mock.Anything, mock.MatchedBy(func(s *domain.DiningSession) bool {
			return s.TableID == tableID && s.TotalAmount == 20.0 && s.ExpiresAt != nil
		})).Return(&domain.DiningSession{ID: uuid.New()}, nil).Once()

		_, err := usecase.CreateSession(context.Background(), input)

		assert.NoError(t, err)
		mockSessionRepo.AssertExpectations(t)
		mockTypeRepo.AssertExpectations(t)
	})
}
