package model

import "errors"

var (
	ErrInvalidID   = errors.New("invalid id")
	ErrNotFound    = errors.New("not found")
	ErrInvalidTask = errors.New("invalid task")
)
