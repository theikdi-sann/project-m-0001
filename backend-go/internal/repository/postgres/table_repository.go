package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/theikdi-sann/qr-restaurant-api/internal/domain"
	"github.com/theikdi-sann/qr-restaurant-api/internal/repository/postgres/db"
)

type tableRepository struct {
	queries *db.Queries
}

func NewTableRepository(queries *db.Queries) domain.TableRepository {
	return &tableRepository{
		queries: queries,
	}
}

func (r *tableRepository) ListAll(ctx context.Context) ([]*domain.DiningTable, error) {
	rows, err := r.queries.ListTables(ctx)
	if err != nil {
		return nil, err
	}

	var tables []*domain.DiningTable
	for _, row := range rows {
		tables = append(tables, &domain.DiningTable{
			ID:          pgToUuid(row.ID),
			TableNumber: row.TableNumber,
			Capacity:    int(row.Capacity),
			Status:      row.Status,
		})
	}
	return tables, nil
}

func (r *tableRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.DiningTable, error) {
	row, err := r.queries.GetTable(ctx, uuidToPg(id))
	if err != nil {
		return nil, err
	}
	return &domain.DiningTable{
		ID:          pgToUuid(row.ID),
		TableNumber: row.TableNumber,
		Capacity:    int(row.Capacity),
		Status:      row.Status,
	}, nil
}
