package cargonaut

import "errors"

var (
	// ErrTokenExists is raised when a token with the same unique constraints
	// already exists.
	ErrTokenExists = errors.New("token exists")
	// ErrTokenNotFound is raised when a token does not exist.
	ErrTokenNotFound = errors.New("token not found")
	// ErrUserExists is raised when a user with the same unique constraints
	// already exists.
	ErrUserExists = errors.New("user exists")
	// ErrUserNotFound is raised when a user does not exist.
	ErrUserNotFound = errors.New("user not found")
	// ErrRatingExists is raised when a rating with the same unique constraints
	// already exists.
	ErrRatingExists = errors.New("rating exists")
	// ErrRatingNotFound is raised when a rating does not exist.
	ErrRatingNotFound = errors.New("rating not found")
	// ErrVehicleExists is raised when a vehicle with the same unique
	// constraints already exists.
	ErrVehicleExists = errors.New("vehicle exists")
	// ErrVehicleNotFound is raised when a vehicle does not exist.
	ErrVehicleNotFound = errors.New("vehicle not found")
	// ErrTripExists is raised when a trip with the same unique constraints
	// already exists.
	ErrTripExists = errors.New("trip exists")
	// ErrTripNotFound is raised when a trip does not exist.
	ErrTripNotFound = errors.New("trip not found")
)
