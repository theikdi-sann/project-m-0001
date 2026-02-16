package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/theikdi-sann/qr-restaurant-api/internal/domain"
)

// Mocks for Order Tests
type MockOrderRepo struct {
	mock.Mock
}

func (m *MockOrderRepo) Create(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	args := m.Called(ctx, order)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Order), args.Error(1)
}
func (m *MockOrderRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Order), args.Error(1)
}
func (m *MockOrderRepo) ListBySessionID(ctx context.Context, sessionID uuid.UUID) ([]*domain.Order, error) {
	return nil, nil
}
func (m *MockOrderRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.OrderStatus) (*domain.Order, error) {
	args := m.Called(ctx, id, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Order), args.Error(1)
}

type MockSessionRepoOrder struct {
	mock.Mock
}

func (m *MockSessionRepoOrder) Create(ctx context.Context, session *domain.DiningSession) (*domain.DiningSession, error) {
	return nil, nil
}
func (m *MockSessionRepoOrder) GetByID(ctx context.Context, id uuid.UUID) (*domain.DiningSession, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.DiningSession), args.Error(1)
}
func (m *MockSessionRepoOrder) GetActiveSessionByTableID(ctx context.Context, tableID uuid.UUID) (*domain.DiningSession, error) {
	return nil, nil
}
func (m *MockSessionRepoOrder) ListExpiredActiveSessions(ctx context.Context) ([]*domain.DiningSession, error) {
	return nil, nil
}
func (m *MockSessionRepoOrder) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.SessionStatus) (*domain.DiningSession, error) {
	return nil, nil
}
func (m *MockSessionRepoOrder) Extend(ctx context.Context, id uuid.UUID, newExpiry time.Time) (*domain.DiningSession, error) {
	return nil, nil
}

type MockMenuItemRepo struct {
	mock.Mock
}

func (m *MockMenuItemRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.MenuItem, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.MenuItem), args.Error(1)
}
func (m *MockMenuItemRepo) ListByCategory(ctx context.Context, categoryID uuid.UUID) ([]*domain.MenuItem, error) {
	return nil, nil
}
func (m *MockMenuItemRepo) ListCategories(ctx context.Context) ([]*domain.MenuCategory, error) {
	return nil, nil
}

func TestCreateOrder(t *testing.T) {
	mockOrderRepo := new(MockOrderRepo)
	mockSessionRepo := new(MockSessionRepoOrder)
	mockMenuRepo := new(MockMenuItemRepo)

	usecase := NewOrderUsecase(mockOrderRepo, mockSessionRepo, mockMenuRepo)

	sessionID := uuid.New()
	menuItemID := uuid.New()

	input := CreateOrderInput{
		DiningSessionID: sessionID,
		Items: []CreateOrderItemInput{
			{MenuItemID: menuItemID, Quantity: 2},
		},
	}

	t.Run("should fail if session is not active", func(t *testing.T) {
		mockSessionRepo.On("GetByID", mock.Anything, sessionID).Return(&domain.DiningSession{
			ID:     sessionID,
			Status: domain.SessionStatusCompleted,
		}, nil).Once()

		_, err := usecase.CreateOrder(context.Background(), input)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrSessionNotActive, err)
	})

	t.Run("should fail if menu item not found", func(t *testing.T) {
		mockSessionRepo.On("GetByID", mock.Anything, sessionID).Return(&domain.DiningSession{
			ID:     sessionID,
			Status: domain.SessionStatusActive,
		}, nil).Once()

		mockMenuRepo.On("GetByID", mock.Anything, menuItemID).Return(nil, errors.New("not found")).Once()

		_, err := usecase.CreateOrder(context.Background(), input)
		assert.Error(t, err)
	})

	t.Run("should create order successfully", func(t *testing.T) {
		mockSessionRepo.On("GetByID", mock.Anything, sessionID).Return(&domain.DiningSession{
			ID:     sessionID,
			Status: domain.SessionStatusActive,
		}, nil).Once()

		price := decimal.NewFromFloat(10.0)
		mockMenuRepo.On("GetByID", mock.Anything, menuItemID).Return(&domain.MenuItem{
			ID:          menuItemID,
			Price:       price,
			IsAvailable: true,
		}, nil).Once()

		expectedTotal := price.Mul(decimal.NewFromInt(2))

		mockOrderRepo.On("Create", mock.Anything, mock.MatchedBy(func(o *domain.Order) bool {
			return o.TotalAmount.Equal(expectedTotal) && len(o.Items) == 1
		})).Return(&domain.Order{ID: uuid.New()}, nil).Once()

		_, err := usecase.CreateOrder(context.Background(), input)
		assert.NoError(t, err)
	})
}

func TestUpdateOrderStatus(t *testing.T) {
	mockOrderRepo := new(MockOrderRepo)
	usecase := NewOrderUsecase(mockOrderRepo, nil, nil) // Other repos not needed for this

	orderID := uuid.New()

	t.Run("should update status successfully for valid transition", func(t *testing.T) {
		// Mock GetByID
		mockOrderRepo.On("GetByID", mock.Anything, orderID).Return(&domain.Order{
			ID:     orderID,
			Status: domain.OrderStatusPending,
		}, nil).Once()

		// Mock UpdateStatus
		mockOrderRepo.On("UpdateStatus", mock.Anything, orderID, domain.OrderStatusPreparing).
			Return(&domain.Order{ID: orderID, Status: domain.OrderStatusPreparing}, nil).Once()

		_, err := usecase.UpdateOrderStatus(context.Background(), orderID, domain.OrderStatusPreparing)
		assert.NoError(t, err)
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("should fail for invalid transition", func(t *testing.T) {
		// Pending -> Served (Invalid, must be Prepared -> Ready first)
		mockOrderRepo.On("GetByID", mock.Anything, orderID).Return(&domain.Order{
			ID:     orderID,
			Status: domain.OrderStatusPending,
		}, nil).Once()

		_, err := usecase.UpdateOrderStatus(context.Background(), orderID, domain.OrderStatusServed)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidStatusTransition, err)
		mockOrderRepo.AssertNotCalled(t, "UpdateStatus")
	})

	t.Run("should allow cancellation from any active state", func(t *testing.T) {
		mockOrderRepo.On("GetByID", mock.Anything, orderID).Return(&domain.Order{
			ID:     orderID,
			Status: domain.OrderStatusPreparing,
		}, nil).Once()

		mockOrderRepo.On("UpdateStatus", mock.Anything, orderID, domain.OrderStatusCancelled).
			Return(&domain.Order{ID: orderID, Status: domain.OrderStatusCancelled}, nil).Once()

		_, err := usecase.UpdateOrderStatus(context.Background(), orderID, domain.OrderStatusCancelled)
		assert.NoError(t, err)
	})
}
