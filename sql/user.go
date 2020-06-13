package sql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"

	"github.com/my-cargonaut/cargonaut"
	_ "github.com/my-cargonaut/cargonaut/sql/migrations" // Migrations
)

var _ cargonaut.UserService = (*UserService)(nil)

const (
	listUsersSQL      = "SELECT id, email, password_hash, display_name, created_at, updated_at FROM user_account ORDER BY updated_at DESC"
	getUserSQL        = "SELECT id, email, password_hash, display_name, created_at, updated_at FROM user_account WHERE id = $1 LIMIT 1"
	getUserByEmailSQL = "SELECT id, email, password_hash, display_name, created_at, updated_at FROM user_account WHERE email = $1 LIMIT 1"
	createUserSQL     = "INSERT INTO user_account (email, password_hash, display_name) VALUES (:email, :password_hash, :display_name)"
	updateUserSQL     = "UPDATE user_account SET email = :email, password_hash = :password_hash, display_name = :display_name, updated_at = :updated_at WHERE id = :id"
	deleteUserSQL     = "DELETE FROM user_account WHERE id = $1"
)

// UserService provides access to the user resource backed by a Postgres SQL
// database.
type UserService struct {
	db *sqlx.DB

	listStmt       *sqlx.Stmt
	getStmt        *sqlx.Stmt
	getByEmailStmt *sqlx.Stmt
	createStmt     *sqlx.NamedStmt
	updateStmt     *sqlx.NamedStmt
	deleteStmt     *sqlx.Stmt
}

// NewUserService returns a new UserService based on top of the provided
// database connection.
func NewUserService(ctx context.Context, db *sqlx.DB) (*UserService, error) {
	s := &UserService{db: db}

	var err error
	if s.listStmt, err = db.PreparexContext(ctx, listUsersSQL); err != nil {
		return nil, fmt.Errorf("prepare list statement: %w", err)
	}
	if s.getStmt, err = db.PreparexContext(ctx, getUserSQL); err != nil {
		return nil, fmt.Errorf("prepare get statement: %w", err)
	}
	if s.getByEmailStmt, err = db.PreparexContext(ctx, getUserByEmailSQL); err != nil {
		return nil, fmt.Errorf("prepare get by email statement: %w", err)
	}
	if s.createStmt, err = db.PrepareNamedContext(ctx, createUserSQL); err != nil {
		return nil, fmt.Errorf("prepare create statement: %w", err)
	}
	if s.updateStmt, err = db.PrepareNamedContext(ctx, updateUserSQL); err != nil {
		return nil, fmt.Errorf("prepare edit statement: %w", err)
	}
	if s.deleteStmt, err = db.PreparexContext(ctx, deleteUserSQL); err != nil {
		return nil, fmt.Errorf("prepare delete statement: %w", err)
	}

	return s, nil
}

// Close the user service and close all prepared statements.
func (s *UserService) Close() error {
	if err := s.listStmt.Close(); err != nil {
		return fmt.Errorf("close list statement: %w", err)
	}
	if err := s.getStmt.Close(); err != nil {
		return fmt.Errorf("close get statement: %w", err)
	}
	if err := s.getByEmailStmt.Close(); err != nil {
		return fmt.Errorf("close get by email statement: %w", err)
	}
	if err := s.createStmt.Close(); err != nil {
		return fmt.Errorf("close create statement: %w", err)
	}
	if err := s.updateStmt.Close(); err != nil {
		return fmt.Errorf("close update statement: %w", err)
	}
	if err := s.deleteStmt.Close(); err != nil {
		return fmt.Errorf("close delete statement: %w", err)
	}
	return nil
}

// ListUsers lists all users.
func (s *UserService) ListUsers(ctx context.Context) ([]*cargonaut.User, error) {
	users := make([]*cargonaut.User, 0)
	if err := s.listStmt.SelectContext(ctx, &users); err != nil {
		return nil, fmt.Errorf("select users from database: %w", err)
	}
	return users, nil
}

// GetUser returns a user identified by his unique ID.
func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (*cargonaut.User, error) {
	user := new(cargonaut.User)
	if err := s.getStmt.GetContext(ctx, user, id); err == sql.ErrNoRows {
		return nil, cargonaut.ErrUserNotFound
	} else if err != nil {
		return nil, fmt.Errorf("get user %q from database: %w", id, err)
	}
	return user, nil
}

// GetUserByEmail returns a user identified by his E-Mail address.
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*cargonaut.User, error) {
	user := new(cargonaut.User)
	if err := s.getByEmailStmt.GetContext(ctx, user, email); err == sql.ErrNoRows {
		return nil, cargonaut.ErrUserNotFound
	} else if err != nil {
		return nil, fmt.Errorf("get user %q from database: %w", email, err)
	}
	return user, nil
}

// CreateUser creates a new user.
func (s *UserService) CreateUser(ctx context.Context, user *cargonaut.User) error {
	if _, err := s.createStmt.ExecContext(ctx, user); isAlreadyExistsError(err) {
		return cargonaut.ErrUserExists
	} else if err != nil {
		return fmt.Errorf("create user in database: %w", err)
	}
	return nil
}

// UpdateUser updates a given user.
func (s *UserService) UpdateUser(ctx context.Context, user *cargonaut.User) error {
	if _, err := s.updateStmt.ExecContext(ctx, user); isAlreadyExistsError(err) {
		return cargonaut.ErrUserExists
	} else if err != nil {
		return fmt.Errorf("update user %q in database: %w", user.ID, err)
	}
	return nil
}

// DeleteUser deletes a user identified by his unique ID.
func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if _, err := s.deleteStmt.ExecContext(ctx, id); err != nil {
		return fmt.Errorf("delete user from database: %w", err)
	}
	return nil
}
