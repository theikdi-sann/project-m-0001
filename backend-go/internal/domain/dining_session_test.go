package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewDiningSession(t *testing.T) {
	tableID := uuid.New()
	sessionTypeID := uuid.New()
	createdBy := uuid.New()
	startTime := time.Now()
	guestCount := 2

	t.Run("should calculate expiry time for buffet session", func(t *testing.T) {
		duration := 90
		buffetType := DiningSessionType{
			ID:              sessionTypeID,
			IsBuffet:        true,
			DurationMinutes: duration,
		}

		session := NewDiningSession(tableID, sessionTypeID, createdBy, guestCount, buffetType, startTime)

		assert.NotNil(t, session.ExpiresAt)
		assert.Equal(t, createdBy, session.CreatedBy)
		expectedExpiry := startTime.Add(time.Duration(duration) * time.Minute)
		assert.True(t, session.ExpiresAt.Equal(expectedExpiry), "Expiry time should be start time + duration")
	})

	t.Run("should calculate total price for buffet session", func(t *testing.T) {
		price := 25.50
		guestCount := 3
		buffetType := DiningSessionType{
			ID:       sessionTypeID,
			IsBuffet: true,
			Price:    price,
		}

		session := NewDiningSession(tableID, sessionTypeID, createdBy, guestCount, buffetType, startTime)

		assert.Equal(t, price, session.PricePerGuest)
		assert.Equal(t, price*float64(guestCount), session.TotalAmount)
	})

	t.Run("should have nil expiry time for non-buffet (à la carte) session", func(t *testing.T) {
		alaCarteType := DiningSessionType{
			ID:       sessionTypeID,
			IsBuffet: false,
		}

		session := NewDiningSession(tableID, sessionTypeID, createdBy, guestCount, alaCarteType, startTime)

		assert.Nil(t, session.ExpiresAt, "À la carte sessions should not have an expiry time")
	})
}
