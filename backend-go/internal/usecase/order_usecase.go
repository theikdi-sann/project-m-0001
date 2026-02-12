package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/theikdi-sann/qr-restaurant-api/internal/domain"
)

type OrderUsecase interface {
	CreateOrder(ctx context.Context, input CreateOrderInput) (*domain.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, newStatus domain.OrderStatus) (*domain.Order, error)
}

type CreateOrderInput struct {
	DiningSessionID uuid.UUID
	Items           []CreateOrderItemInput
}

type CreateOrderItemInput struct {
	MenuItemID uuid.UUID
	Quantity   int
	Notes      string
}

type orderUsecase struct {
	orderRepo   domain.OrderRepository
	sessionRepo domain.DiningSessionRepository
	menuRepo    domain.MenuItemRepository
}

func NewOrderUsecase(o domain.OrderRepository, s domain.DiningSessionRepository, m domain.MenuItemRepository) OrderUsecase {
	return &orderUsecase{
		orderRepo:   o,
		sessionRepo: s,
		menuRepo:    m,
	}
}

func (u *orderUsecase) CreateOrder(ctx context.Context, input CreateOrderInput) (*domain.Order, error) {
	// 1. Validate Session
	session, err := u.sessionRepo.GetByID(ctx, input.DiningSessionID)
	if err != nil {
		return nil, err
	}
	if session.Status != domain.SessionStatusActive {
		return nil, domain.ErrSessionNotActive
	}
	if session.ExpiresAt != nil && time.Now().After(*session.ExpiresAt) {
		return nil, domain.ErrSessionExpired
	}

	// 2. Process Items & Calculate Total
	var orderItems []domain.OrderItem
	totalAmount := decimal.Zero

	for _, itemInput := range input.Items {
		menuItem, err := u.menuRepo.GetByID(ctx, itemInput.MenuItemID)
		if err != nil {
			return nil, err // Item not found or DB error
		}
		if !menuItem.IsAvailable {
			return nil, domain.ErrNotFound // Or specific ErrItemUnavailable
		}

		// Snapshot price
		unitPrice := menuItem.Price
		lineTotal := unitPrice.Mul(decimal.NewFromInt(int64(itemInput.Quantity)))
		totalAmount = totalAmount.Add(lineTotal)

		orderItems = append(orderItems, domain.OrderItem{
			ID:         uuid.New(),
			MenuItemID: itemInput.MenuItemID,
			Quantity:   itemInput.Quantity,
			UnitPrice:  unitPrice,
			Notes:      itemInput.Notes,
		})
	}

	// 3. Create Order Domain Object
	order := &domain.Order{
		ID:              uuid.New(),
		DiningSessionID: input.DiningSessionID,
		Status:          domain.OrderStatusPending,
		TotalAmount:     totalAmount,
		Items:           orderItems,
		CreatedAt:       time.Now(),
	}

	// 4. Save
	return u.orderRepo.Create(ctx, order)
}

func (u *orderUsecase) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, newStatus domain.OrderStatus) (*domain.Order, error) {
	order, err := u.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, domain.ErrNotFound
	}

	if !isValidTransition(order.Status, newStatus) {
		return nil, domain.ErrInvalidStatusTransition
	}

	return u.orderRepo.UpdateStatus(ctx, orderID, newStatus)
}

func isValidTransition(current, next domain.OrderStatus) bool {
	if current == next {
		return true // No change
	}
	if next == domain.OrderStatusCancelled {
		return current != domain.OrderStatusServed && current != domain.OrderStatusCancelled
	}

	switch current {
	case domain.OrderStatusPending:
		return next == domain.OrderStatusPreparing
	case domain.OrderStatusPreparing:
		return next == domain.OrderStatusReady
	case domain.OrderStatusReady:
		return next == domain.OrderStatusServed
	default:
		return false
	}
}
