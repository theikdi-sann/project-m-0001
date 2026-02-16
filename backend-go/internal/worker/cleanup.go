package worker

import (
	"context"
	"log"
	"time"

	"github.com/theikdi-sann/qr-restaurant-api/internal/domain"
)

type SessionCleanupWorker struct {
	sessionRepo domain.DiningSessionRepository
}

func NewSessionCleanupWorker(repo domain.DiningSessionRepository) *SessionCleanupWorker {
	return &SessionCleanupWorker{
		sessionRepo: repo,
	}
}

func (w *SessionCleanupWorker) Start(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				w.cleanup()
			}
		}
	}()
}

func (w *SessionCleanupWorker) cleanup() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sessions, err := w.sessionRepo.ListExpiredActiveSessions(ctx)
	if err != nil {
		log.Printf("Worker Error: Failed to list expired sessions: %v", err)
		return
	}

	for _, s := range sessions {
		log.Printf("Worker: Expiring session %s", s.ID)
		_, err := w.sessionRepo.UpdateStatus(ctx, s.ID, domain.SessionStatusExpired)
		if err != nil {
			log.Printf("Worker Error: Failed to expire session %s: %v", s.ID, err)
		}
	}
}
