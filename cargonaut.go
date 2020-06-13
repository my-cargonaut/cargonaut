package cargonaut

import (
	"context"
	"time"

	uuid "github.com/satori/go.uuid"
)

// User represents a user identity.
type User struct {
	ID          uuid.UUID `json:"id" db:"id" sql:"type:uuid"`
	Email       string    `json:"email" db:"email"`
	Password    string    `json:"-" db:"password_hash"`
	DisplayName string    `json:"display_name" db:"display_name"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

// UserService provides access to the user resource.
type UserService interface {
	// ListUsers lists all users.
	ListUsers(context.Context) ([]*User, error)
	// GetUser returns a user identified by his unique ID.
	GetUser(ctx context.Context, id uuid.UUID) (*User, error)
	// GetUserByEmail returns a user identified by his E-Mail address.
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	// CreateUser creates a new user.
	CreateUser(context.Context, *User) error
	// UpdateUser updates a given user.
	UpdateUser(context.Context, *User) error
	// DeleteUser deletes a user identified by his unique ID.
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
