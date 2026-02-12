package domain

import "errors"

var (
	ErrTableOccupied    = errors.New("table is already occupied")
	ErrNotFound         = errors.New("not found")
	ErrSessionNotActive = errors.New("session is not active")
	ErrSessionExpired   = errors.New("session has expired")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
)
