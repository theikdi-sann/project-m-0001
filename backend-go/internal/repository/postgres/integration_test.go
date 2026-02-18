package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/theikdi-sann/qr-restaurant-api/internal/domain"
	"github.com/theikdi-sann/qr-restaurant-api/internal/repository/postgres/db"
)

const dbSource = "postgresql://postgres:postgres@127.0.0.1:54328/postgres"

func setupTestDB(t *testing.T) *pgxpool.Pool {
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, dbSource)
	require.NoError(t, err)

	err = conn.Ping(ctx)
	require.NoError(t, err)

	return conn
}

func TestDiningSessionRepositoryIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	conn := setupTestDB(t)
	defer conn.Close()

	repo := NewDiningSessionRepository(conn)
	ctx := context.Background()

	// 1. Seed Dependencies
	tableID := uuid.New()
	sessionTypeID := uuid.New()
	userID := uuid.New()

	// Create User
	_, err := conn.Exec(ctx, "INSERT INTO users (id, email, role) VALUES ($1, $2, $3)", userID, "test@example.com", "manager")
	require.NoError(t, err)
	defer conn.Exec(ctx, "DELETE FROM users WHERE id = $1", userID)

	// Create Table
	_, err = conn.Exec(ctx, "INSERT INTO dining_tables (id, table_number) VALUES ($1, $2)", tableID, "T-99")
	require.NoError(t, err)
	defer conn.Exec(ctx, "DELETE FROM dining_tables WHERE id = $1", tableID)

	// Create Session Type
	_, err = conn.Exec(ctx, "INSERT INTO dining_session_types (id, name, is_buffet, price, duration_minutes) VALUES ($1, $2, $3, $4, $5)", sessionTypeID, "Integration Buffet", true, 29.99, 90)
	require.NoError(t, err)
	defer conn.Exec(ctx, "DELETE FROM dining_session_types WHERE id = $1", sessionTypeID)

	t.Run("Create and Get DiningSession", func(t *testing.T) {
		// Prepare Session Domain Object
		startTime := time.Now().Round(time.Microsecond) // Postgres resolution
		expiry := startTime.Add(90 * time.Minute)

		session := &domain.DiningSession{
			TableID:       tableID,
			SessionTypeID: sessionTypeID,
			CreatedBy:     userID,
			GuestCount:    4,
			Status:        domain.SessionStatusActive,
			StartTime:     startTime,
			ExpiresAt:     &expiry,
			PricePerGuest: decimal.NewFromFloat(29.99),
			TotalAmount:   decimal.NewFromFloat(119.96),
		}

		// Test Create
		createdSession, err := repo.Create(ctx, session)
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, createdSession.ID)
		assert.Equal(t, session.TableID, createdSession.TableID)
		// Compare decimals using Equal to avoid precision issues
		assert.True(t, session.TotalAmount.Equal(createdSession.TotalAmount))

		defer conn.Exec(ctx, "DELETE FROM dining_sessions WHERE id = $1", createdSession.ID)

		// Test GetByID
		fetchedSession, err := repo.GetByID(ctx, createdSession.ID)
		require.NoError(t, err)
		assert.Equal(t, createdSession.ID, fetchedSession.ID)
		assert.Equal(t, createdSession.GuestCount, fetchedSession.GuestCount)
	})

	t.Run("GetActiveSessionByTableID", func(t *testing.T) {
		// 1. Create a session first
		startTime := time.Now()
		session := &domain.DiningSession{
			TableID:       tableID,
			SessionTypeID: sessionTypeID,
			CreatedBy:     userID,
			GuestCount:    2,
			Status:        domain.SessionStatusActive,
			StartTime:     startTime,
			PricePerGuest: decimal.NewFromFloat(29.99),
			TotalAmount:   decimal.NewFromFloat(59.98),
		}
		createdSession, err := repo.Create(ctx, session)
		require.NoError(t, err)
		defer conn.Exec(ctx, "DELETE FROM dining_sessions WHERE id = $1", createdSession.ID)

		// 2. Fetch Active
		activeSession, err := repo.GetActiveSessionByTableID(ctx, tableID)
		require.NoError(t, err)
		assert.NotNil(t, activeSession)
		assert.Equal(t, createdSession.ID, activeSession.ID)

		// 3. Update Status to Completed
		updatedSession, err := repo.UpdateStatus(ctx, activeSession.ID, domain.SessionStatusCompleted)
		require.NoError(t, err)
		assert.Equal(t, domain.SessionStatusCompleted, updatedSession.Status)

		// 4. Fetch Active again (Should be nil)
		activeSessionAfter, err := repo.GetActiveSessionByTableID(ctx, tableID)
		require.NoError(t, err)
		assert.Nil(t, activeSessionAfter, "Should return nil if no active session exists")
	})

	t.Run("Take Away Session Lifecycle", func(t *testing.T) {
		// 1. Create Take Away Type (Duration NULL)
		takeAwayTypeID := uuid.New()
		_, err := conn.Exec(ctx, "INSERT INTO dining_session_types (id, name, is_buffet, price, duration_minutes) VALUES ($1, $2, $3, $4, NULL)", takeAwayTypeID, "Integration TakeAway", false, 0.00)
		require.NoError(t, err)
		defer conn.Exec(ctx, "DELETE FROM dining_session_types WHERE id = $1", takeAwayTypeID)

		// 2. Create Session (Table Nil)
		startTime := time.Now()
		session := &domain.DiningSession{
			TableID:       uuid.Nil, // No Table
			SessionTypeID: takeAwayTypeID,
			CreatedBy:     userID,
			GuestCount:    1,
			Status:        domain.SessionStatusActive,
			StartTime:     startTime,
			ExpiresAt:     nil, // Unlimited
			PricePerGuest: decimal.Zero,
			TotalAmount:   decimal.Zero,
		}

		createdSession, err := repo.Create(ctx, session)
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, createdSession.ID)
		assert.Equal(t, uuid.Nil, createdSession.TableID) // Should be Nil
		assert.Nil(t, createdSession.ExpiresAt)

		defer conn.Exec(ctx, "DELETE FROM dining_sessions WHERE id = $1", createdSession.ID)
	})

	t.Run("Menu Availability Lifecycle", func(t *testing.T) {
		repo := NewMenuItemRepository(db.New(conn))

		// 1. Create Category
		catID := uuid.New()
		_, err := conn.Exec(ctx, "INSERT INTO menu_categories (id, name, sort_order) VALUES ($1, $2, $3)", catID, "Test Cat", 1)
		require.NoError(t, err)
		defer conn.Exec(ctx, "DELETE FROM menu_categories WHERE id = $1", catID)

		// 2. Create Item (Available)
		itemID := uuid.New()
		_, err = conn.Exec(ctx, "INSERT INTO menu_items (id, category_id, name, price, is_available) VALUES ($1, $2, $3, $4, $5)", itemID, catID, "Test Item", 10.0, true)
		require.NoError(t, err)
		defer conn.Exec(ctx, "DELETE FROM menu_items WHERE id = $1", itemID)

		// 3. Verify ListByCategory (Should find it)
		items, err := repo.ListByCategory(ctx, catID)
		require.NoError(t, err)
		assert.Len(t, items, 1)

		// 4. Update to Unavailable
		updated, err := repo.UpdateAvailability(ctx, itemID, false)
		require.NoError(t, err)
		assert.False(t, updated.IsAvailable)

		// 5. Verify ListByCategory (Should NOT find it)
		items, err = repo.ListByCategory(ctx, catID)
		require.NoError(t, err)
		assert.Len(t, items, 0)

		// 6. Verify ListAllByCategory (Should find it)
		itemsAll, err := repo.ListAllByCategory(ctx, catID)
		require.NoError(t, err)
		assert.Len(t, itemsAll, 1)
		assert.False(t, itemsAll[0].IsAvailable)
	})
}
