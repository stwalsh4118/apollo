package repository

import "errors"

// Sentinel errors for write operations.
var (
	ErrNotFound       = errors.New("not found")
	ErrDuplicate      = errors.New("duplicate entry")
	ErrFKViolation    = errors.New("foreign key violation")
	ErrCheckViolation = errors.New("check constraint violation")
)
