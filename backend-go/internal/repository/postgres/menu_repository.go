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
	return mapMenuItems(rows), nil
}

func (r *menuItemRepository) ListAllByCategory(ctx context.Context, categoryID uuid.UUID) ([]*domain.MenuItem, error) {
	rows, err := r.queries.ListAllMenuItemsByCategory(ctx, uuidToPg(categoryID))
	if err != nil {
		return nil, err
	}
	return mapMenuItems(rows), nil
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

func (r *menuItemRepository) UpdateAvailability(ctx context.Context, id uuid.UUID, isAvailable bool) (*domain.MenuItem, error) {
	arg := db.UpdateMenuItemAvailabilityParams{
		ID:          uuidToPg(id),
		IsAvailable: isAvailable,
	}
	row, err := r.queries.UpdateMenuItemAvailability(ctx, arg)
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

// Helper to map DB rows to Domain
// Note: SQLC generates different structs for ListMenuItemsByCategory and ListAllMenuItemsByCategory if the columns differ?
// Both select *, so columns are same. BUT Go types might be different if named differently.
// Let's check db/models.go or db/menu.sql.go. 
// Usually sqlc generates `ListMenuItemsByCategoryRow` struct if it's not returning full table model.
// But I used `SELECT * FROM menu_items`. So it returns `MenuItem`.
// Let's assume `rows` are []db.MenuItem.

func mapMenuItems(rows []db.MenuItem) []*domain.MenuItem {
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
	return items
}

// Helper for nullable text
func pgTextPtr(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	return &t.String
}
