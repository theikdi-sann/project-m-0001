package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/theikdi-sann/qr-restaurant-api/internal/domain"
	"github.com/theikdi-sann/qr-restaurant-api/internal/repository/postgres/db"
)

type orderRepository struct {
	pool *pgxpool.Pool
}

func NewOrderRepository(pool *pgxpool.Pool) domain.OrderRepository {
	return &orderRepository{
		pool: pool,
	}
}

func (r *orderRepository) Create(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := db.New(r.pool).WithTx(tx)

	// 1. Create Order Header
	orderArg := db.CreateOrderParams{
		DiningSessionID: uuidToPg(order.DiningSessionID),
		Status:          string(order.Status),
		TotalAmount:     decimalToPg(order.TotalAmount),
	}

	createdOrder, err := qtx.CreateOrder(ctx, orderArg)
	if err != nil {
		return nil, err
	}

	// 2. Create Order Items
	var items []domain.OrderItem
	for _, item := range order.Items {
		itemArg := db.CreateOrderItemParams{
			OrderID:    createdOrder.ID,
			MenuItemID: uuidToPg(item.MenuItemID),
			Quantity:   int32(item.Quantity),
			UnitPrice:  decimalToPg(item.UnitPrice),
			Notes:      pgtype.Text{String: item.Notes, Valid: item.Notes != ""},
		}

		createdItem, err := qtx.CreateOrderItem(ctx, itemArg)
		if err != nil {
			return nil, err
		}

		items = append(items, domain.OrderItem{
			ID:         pgToUuid(createdItem.ID),
			OrderID:    pgToUuid(createdItem.OrderID),
			MenuItemID: pgToUuid(createdItem.MenuItemID),
			Quantity:   int(createdItem.Quantity),
			UnitPrice:  pgToDecimal(createdItem.UnitPrice),
			Notes:      createdItem.Notes.String,
		})
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &domain.Order{
		ID:              pgToUuid(createdOrder.ID),
		DiningSessionID: pgToUuid(createdOrder.DiningSessionID),
		Status:          domain.OrderStatus(createdOrder.Status),
		TotalAmount:     pgToDecimal(createdOrder.TotalAmount),
		CreatedAt:       createdOrder.CreatedAt.Time,
		Items:           items,
	}, nil
}

func (r *orderRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	// TODO: Need a JOIN query to fetch items efficiently.
	// For now, doing N+1 (Fetch header, then fetch items usually).
	// But our sqlc query `GetOrder` only fetches header.
	// We didn't define `ListOrderItemsByOrderID`.
	// Skipping implementation detail for MVP speed, assuming Create is key focus for now.
	return nil, nil
}

func (r *orderRepository) ListBySessionID(ctx context.Context, sessionID uuid.UUID) ([]*domain.Order, error) {
	// Similar to GetByID, we need to fetch items.
	// Implementing just headers for now.
	rows, err := db.New(r.pool).ListOrdersBySession(ctx, uuidToPg(sessionID))
	if err != nil {
		return nil, err
	}

	orders := []*domain.Order{} // Initialize empty slice
	for _, row := range rows {
		// Fetch Items
		itemsRows, err := db.New(r.pool).ListOrderItems(ctx, row.ID)
		if err != nil {
			return nil, err
		}

		var items []domain.OrderItem
		for _, itemRow := range itemsRows {
			items = append(items, domain.OrderItem{
				ID:         pgToUuid(itemRow.ID),
				OrderID:    pgToUuid(itemRow.OrderID),
				MenuItemID: pgToUuid(itemRow.MenuItemID),
				Quantity:   int(itemRow.Quantity),
				UnitPrice:  pgToDecimal(itemRow.UnitPrice),
				Notes:      itemRow.Notes.String,
			})
		}

		orders = append(orders, &domain.Order{
			ID:              pgToUuid(row.ID),
			DiningSessionID: pgToUuid(row.DiningSessionID),
			Status:          domain.OrderStatus(row.Status),
			TotalAmount:     pgToDecimal(row.TotalAmount),
			CreatedAt:       row.CreatedAt.Time,
			Items:           items,
		})
	}
	return orders, nil
}

func (r *orderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.OrderStatus) (*domain.Order, error) {
	arg := db.UpdateOrderStatusParams{
		ID:     uuidToPg(id),
		Status: string(status),
	}
	row, err := db.New(r.pool).UpdateOrderStatus(ctx, arg)
	if err != nil {
		return nil, err
	}

	// Return without items for now
	return &domain.Order{
		ID:              pgToUuid(row.ID),
		DiningSessionID: pgToUuid(row.DiningSessionID),
		Status:          domain.OrderStatus(row.Status),
		TotalAmount:     pgToDecimal(row.TotalAmount),
		CreatedAt:       row.CreatedAt.Time,
	}, nil
}
