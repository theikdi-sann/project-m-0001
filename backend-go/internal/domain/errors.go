package domain

import "errors"

var (
	ErrTableOccupied = errors.New("table is already occupied")
	ErrNotFound      = errors.New("not found")
)
