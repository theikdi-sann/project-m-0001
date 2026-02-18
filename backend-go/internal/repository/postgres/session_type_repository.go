package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/theikdi-sann/qr-restaurant-api/internal/domain"
	"github.com/theikdi-sann/qr-restaurant-api/internal/repository/postgres/db"
)

type sessionTypeRepository struct {
	queries *db.Queries
}

func NewSessionTypeRepository(queries *db.Queries) domain.DiningSessionTypeRepository {
	return &sessionTypeRepository{
		queries: queries,
	}
}

func (r *sessionTypeRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.DiningSessionType, error) {
	row, err := r.queries.GetSessionType(ctx, uuidToPg(id))
	if err != nil {
		return nil, err
	}

	return &domain.DiningSessionType{
		ID:              pgToUuid(row.ID),
		Name:            row.Name,
		IsBuffet:        row.IsBuffet,
		Price:           pgToDecimal(row.Price),
		DurationMinutes: int(row.DurationMinutes.Int32),
	}, nil
}

func (r *sessionTypeRepository) ListAll(ctx context.Context) ([]*domain.DiningSessionType, error) {
	rows, err := r.queries.ListSessionTypes(ctx)
	if err != nil {
		return nil, err
	}

	var types []*domain.DiningSessionType
	for _, row := range rows {
		types = append(types, &domain.DiningSessionType{
			ID:              pgToUuid(row.ID),
			Name:            row.Name,
			IsBuffet:        row.IsBuffet,
			Price:           pgToDecimal(row.Price),
			DurationMinutes: int(row.DurationMinutes.Int32),
		})
	}
	return types, nil
}
