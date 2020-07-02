package cargonaut

import (
	"context"
	"time"

	uuid "github.com/satori/go.uuid"
)

// Rating is a rating given by a user (Author) to another user. Users can't rate
// themselves.
type Rating struct {
	ID        uuid.UUID `json:"id" db:"id" sql:"type:uuid"`
	UserID    uuid.UUID `json:"user_id" db:"user_id" sql:"type:uuid"`
	AuthorID  uuid.UUID `json:"author_id" db:"author_id" sql:"type:uuid"`
	Comment   string    `json:"comment" db:"comment"`
	Value     float32   `json:"value" db:"value"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Token represents an authentication token.
type Token struct {
	ID        uuid.UUID `db:"id" sql:"type:uuid"`
	UserID    uuid.UUID `db:"user_id" sql:"type:uuid"`
	Token     string    `db:"-"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
}

// Trip is a trip from one location to another one. If no RiderID is assigned,
// the trip is still available.
type Trip struct {
	ID          uuid.UUID  `json:"id" db:"id" sql:"type:uuid"`
	UserID      uuid.UUID  `json:"user_id" db:"user_id" sql:"type:uuid"`
	VehicleID   uuid.UUID  `json:"vehicle_id" db:"vehicle_id" sql:"type:uuid"`
	RiderID     *uuid.UUID `json:"rider_id" db:"rider_id" sql:"type:uuid"`
	Start       string     `json:"start" db:"start"`
	Destination string     `json:"destination" db:"destination"`
	Price       float32    `json:"price" db:"price"`
	Depature    time.Time  `json:"depature" db:"depature"`
	Arrival     time.Time  `json:"arrival" db:"arrival"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// User represents a user identity.
type User struct {
	ID          uuid.UUID `json:"id" db:"id" sql:"type:uuid"`
	Email       string    `json:"email" db:"email"`
	Password    string    `json:"-" db:"password_hash"`
	DisplayName string    `json:"display_name" db:"display_name"`
	Birthday    time.Time `json:"birthday" db:"birthday"`
	Avatar      string    `json:"-" db:"avatar"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Vehicle is a vehicle belonging to a user.
type Vehicle struct {
	ID                uuid.UUID `json:"id" db:"id" sql:"type:uuid"`
	UserID            uuid.UUID `json:"user_id" db:"user_id" sql:"type:uuid"`
	Brand             string    `json:"brand" db:"brand"`
	Model             string    `json:"model" db:"model"`
	Passengers        uint8     `json:"passengers" db:"passengers"`
	LoadingAreaLength float32   `json:"loading_area_length" db:"loading_area_length"`
	LoadingAreaWidth  float32   `json:"loading_area_width" db:"loading_area_width"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

// TokenBlacklist provides methods for blacklisting authentication tokens.
type TokenBlacklist interface {
	// IsTokenBlacklisted retrieves a token by its unique token ID. If the token
	// is on the blacklist, true is returned.
	IsTokenBlacklisted(ctx context.Context, tokenID uuid.UUID) (bool, error)
	// BlacklistToken blacklists one or more tokens by putting them onto the
	// token blacklist. The tokens are identified by their unique token ID.
	BlacklistToken(ctx context.Context, tokens ...*Token) error
}

// TripRepository provides access to the trip resource.
type TripRepository interface {
	// ListTrips lists all trips.
	ListTrips(context.Context) ([]*Trip, error)
	// GetTrip returns a trip identified by its unique ID.
	GetTrip(ctx context.Context, id uuid.UUID) (*Trip, error)
	// CreateTrip creates a new trip.
	CreateTrip(context.Context, *Trip) error
	// UpdateTrip updates a given trip.
	UpdateTrip(context.Context, *Trip) error
	// DeleteTrip deletes a trip identified by his unique ID.
	DeleteTrip(ctx context.Context, id uuid.UUID) error
}

// UserRepository provides access to the user resource.
type UserRepository interface {
	// ListUsers lists all users.
	ListUsers(context.Context) ([]*User, error)
	// GetUser returns a user identified by his unique ID.
	GetUser(ctx context.Context, userID uuid.UUID) (*User, error)
	// GetUserByEmail returns a user identified by his E-Mail address.
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	// CreateUser creates a new user.
	CreateUser(context.Context, *User) error
	// UpdateUser updates a given user.
	UpdateUser(context.Context, *User) error
	// DeleteUser deletes a user identified by his unique ID.
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	// ListTokens lists all authentication tokens for the user identified by his
	// unique ID.
	ListTokens(ctx context.Context, userID uuid.UUID) ([]*Token, error)
	// CreateToken creates an authentication token for the user identified by
	// the tokens unique user ID.
	CreateToken(ctx context.Context, token *Token) error
	// DeleteToken deletes an users authentication token. Token and user are
	// identified by their unique IDs.
	DeleteToken(ctx context.Context, userID, tokenID uuid.UUID) error
	// ListRatings lists all ratings for the user identified by his unique ID.
	ListRatings(ctx context.Context, userID uuid.UUID) ([]*Rating, error)
	// CreateRating creates a new rating.
	CreateRating(context.Context, *Rating) error
	// ListVehicles lists all vehicles for the user identified by his unique ID.
	ListVehicles(ctx context.Context, userID uuid.UUID) ([]*Vehicle, error)
}

// VehicleRepository provides access to the vehicle resource.
type VehicleRepository interface {
	// ListVehicles lists all vehicles.
	ListVehicles(context.Context) ([]*Vehicle, error)
	// GetVehicle returns a vehicle identified by its unique ID.
	GetVehicle(ctx context.Context, id uuid.UUID) (*Vehicle, error)
	// CreateVehicle creates a new vehicle.
	CreateVehicle(context.Context, *Vehicle) error
	// UpdateVehicle updates a given vehicle.
	UpdateVehicle(context.Context, *Vehicle) error
	// DeleteVehicle deletes a vehicle identified by his unique ID.
	DeleteVehicle(ctx context.Context, id uuid.UUID) error
}
