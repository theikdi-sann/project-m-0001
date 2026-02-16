package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/theikdi-sann/qr-restaurant-api/internal/domain"
	"github.com/theikdi-sann/qr-restaurant-api/internal/repository/postgres/db"
)

type menuItemRepository struct {
	queries *db.Queries
}

func NewMenuItemRepository(queries *db.Queries) domain.MenuItemRepository {
	return &menuItemRepository{
		queries: queries,
	}
}

func (r *menuItemRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.MenuItem, error) {
	row, err := r.queries.GetMenuItem(ctx, uuidToPg(id))
	if err != nil {
		return nil, err
	}

	return &domain.MenuItem{
		ID:          pgToUuid(row.ID),
		CategoryID:  pgToUuid(row.CategoryID),
		Name:        row.Name,
		Description: pgTextPtr(row.Description),
		Price:       pgToDecimal(row.Price),
		IsAvailable: row.IsAvailable,
	}, nil
}

func (r *menuItemRepository) ListByCategory(ctx context.Context, categoryID uuid.UUID) ([]*domain.MenuItem, error) {
	rows, err := r.queries.ListMenuItemsByCategory(ctx, uuidToPg(categoryID))
	if err != nil {
		return nil, err
	}

	var items []*domain.MenuItem
	for _, row := range rows {
		items = append(items, &domain.MenuItem{
			ID:          pgToUuid(row.ID),
			CategoryID:  pgToUuid(row.CategoryID),
			Name:        row.Name,
			Description: pgTextPtr(row.Description),
			Price:       pgToDecimal(row.Price),
			IsAvailable: row.IsAvailable,
		})
	}
	return items, nil
}

func (r *menuItemRepository) ListCategories(ctx context.Context) ([]*domain.MenuCategory, error) {
	rows, err := r.queries.ListMenuCategories(ctx)
	if err != nil {
		return nil, err
	}

	var categories []*domain.MenuCategory
	for _, row := range rows {
		categories = append(categories, &domain.MenuCategory{
			ID:        pgToUuid(row.ID),
			Name:      row.Name,
			SortOrder: int(row.SortOrder),
		})
	}
	return categories, nil
}

// Helper for nullable text
func pgTextPtr(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	return &t.String
}
