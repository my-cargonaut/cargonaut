package cargonaut

import "errors"

var (
	// ErrUserExists is raised when a user with the same unique constraints
	// already exists.
	ErrUserExists = errors.New("user exists")
	// ErrUserNotFound is raised when a user does not exist.
	ErrUserNotFound = errors.New("user not found")
)
