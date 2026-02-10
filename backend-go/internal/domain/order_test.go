package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/theikdi-sann/qr-restaurant-api/internal/domain"
)

func TestValidateOrder(t *testing.T) {
	t.Run("should allow order when session is active and not expired", func(t *testing.T) {
		session := &domain.DiningSession{
			ID:     uuid.New(),
			Status: domain.SessionStatusActive,
		}
		order := &domain.Order{
			ID: uuid.New(),
		}

		err := domain.ValidateOrder(session, order)
		assert.NoError(t, err)
	})

	t.Run("should reject order when session is not active", func(t *testing.T) {
		session := &domain.DiningSession{
			ID:     uuid.New(),
			Status: domain.SessionStatusCompleted,
		}
		order := &domain.Order{
			ID: uuid.New(),
		}

		err := domain.ValidateOrder(session, order)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrSessionNotActive, err)
	})

	t.Run("should reject order when session is expired (Buffet)", func(t *testing.T) {
		expiryTime := time.Now().Add(-10 * time.Minute) // Expired 10 mins ago
		session := &domain.DiningSession{
			ID:        uuid.New(),
			Status:    domain.SessionStatusActive,
			ExpiresAt: &expiryTime,
		}
		order := &domain.Order{
			ID: uuid.New(),
		}

		err := domain.ValidateOrder(session, order)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrSessionExpired, err)
	})

	t.Run("should allow order when session has no expiry (A La Carte)", func(t *testing.T) {
		session := &domain.DiningSession{
			ID:        uuid.New(),
			Status:    domain.SessionStatusActive,
			ExpiresAt: nil, // No expiry
		}
		order := &domain.Order{
			ID: uuid.New(),
		}

		err := domain.ValidateOrder(session, order)
		assert.NoError(t, err)
	})
}
