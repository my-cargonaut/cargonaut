package sql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"

	"github.com/my-cargonaut/cargonaut"
	_ "github.com/my-cargonaut/cargonaut/internal/sql/migrations" // Migrations
)

var _ cargonaut.UserRepository = (*UserRepository)(nil)

const (
	listUsersSQL        = "SELECT id, email, password_hash, display_name, birthday, avatar, created_at, updated_at FROM user_account ORDER BY updated_at DESC"
	getUserSQL          = "SELECT id, email, password_hash, display_name, birthday, avatar, created_at, updated_at FROM user_account WHERE id = $1 LIMIT 1"
	getUserByEmailSQL   = "SELECT id, email, password_hash, display_name, birthday, avatar, created_at, updated_at FROM user_account WHERE email = $1 LIMIT 1"
	createUserSQL       = "INSERT INTO user_account (email, password_hash, display_name, birthday, avatar) VALUES (:email, :password_hash, :display_name, :birthday, :avatar)"
	updateUserSQL       = "UPDATE user_account SET email = :email, password_hash = :password_hash, display_name = :display_name, birthday = :birthday, avatar = :avatar, updated_at = :updated_at WHERE id = :id"
	deleteUserSQL       = "DELETE FROM user_account WHERE id = $1"
	listTokensSQL       = "SELECT id, user_id, expires_at, created_at FROM user_token WHERE user_id = $1"
	createTokenSQL      = "INSERT INTO user_token (id, user_id, expires_at) VALUES (:id, :user_id, :expires_at)"
	deleteTokenSQL      = "DELETE FROM user_token WHERE user_id = $1 AND id = $2"
	listRatingsSQL      = "SELECT id, user_id, author_id, trip_id, comment, value, created_at FROM rating WHERE user_id = $1"
	listUserVehiclesSQL = "SELECT id, user_id, brand, model, passengers, loading_area_length, loading_area_width, created_at, updated_at FROM vehicle WHERE user_id = $1 ORDER BY updated_at DESC"
)

// UserRepository provides access to the user resource backed by a Postgres SQL
// database.
type UserRepository struct {
	db *sqlx.DB

	listUsersStmt        *sqlx.Stmt
	getUserStmt          *sqlx.Stmt
	getByEmailUserStmt   *sqlx.Stmt
	createUserStmt       *sqlx.NamedStmt
	updateUserStmt       *sqlx.NamedStmt
	deleteUserStmt       *sqlx.Stmt
	listTokensStmt       *sqlx.Stmt
	createTokenStmt      *sqlx.NamedStmt
	deleteTokenStmt      *sqlx.Stmt
	listRatingsStmt      *sqlx.Stmt
	listUserVehiclesStmt *sqlx.Stmt
}

// NewUserRepository returns a new UserRepository based on top of the provided
// database connection.
func NewUserRepository(ctx context.Context, db *sqlx.DB) (*UserRepository, error) {
	s := &UserRepository{db: db}

	var err error
	if s.listUsersStmt, err = db.PreparexContext(ctx, listUsersSQL); err != nil {
		return nil, fmt.Errorf("prepare list users statement: %w", err)
	}
	if s.getUserStmt, err = db.PreparexContext(ctx, getUserSQL); err != nil {
		return nil, fmt.Errorf("prepare get user statement: %w", err)
	}
	if s.getByEmailUserStmt, err = db.PreparexContext(ctx, getUserByEmailSQL); err != nil {
		return nil, fmt.Errorf("prepare get user by email statement: %w", err)
	}
	if s.createUserStmt, err = db.PrepareNamedContext(ctx, createUserSQL); err != nil {
		return nil, fmt.Errorf("prepare create user statement: %w", err)
	}
	if s.updateUserStmt, err = db.PrepareNamedContext(ctx, updateUserSQL); err != nil {
		return nil, fmt.Errorf("prepare update user statement: %w", err)
	}
	if s.deleteUserStmt, err = db.PreparexContext(ctx, deleteUserSQL); err != nil {
		return nil, fmt.Errorf("prepare delete user statement: %w", err)
	}
	if s.listTokensStmt, err = db.PreparexContext(ctx, listTokensSQL); err != nil {
		return nil, fmt.Errorf("prepare list user tokens statement: %w", err)
	}
	if s.createTokenStmt, err = db.PrepareNamedContext(ctx, createTokenSQL); err != nil {
		return nil, fmt.Errorf("prepare create user token statement: %w", err)
	}
	if s.deleteTokenStmt, err = db.PreparexContext(ctx, deleteTokenSQL); err != nil {
		return nil, fmt.Errorf("prepare delete user token statement: %w", err)
	}
	if s.listRatingsStmt, err = db.PreparexContext(ctx, listRatingsSQL); err != nil {
		return nil, fmt.Errorf("prepare list user ratings statement: %w", err)
	}
	if s.listUserVehiclesStmt, err = db.PreparexContext(ctx, listUserVehiclesSQL); err != nil {
		return nil, fmt.Errorf("prepare list user vehicles statement: %w", err)
	}

	return s, nil
}

// Close all prepared statements.
func (s *UserRepository) Close() error {
	if err := s.listUsersStmt.Close(); err != nil {
		return fmt.Errorf("close list user statement: %w", err)
	}
	if err := s.getUserStmt.Close(); err != nil {
		return fmt.Errorf("close get user statement: %w", err)
	}
	if err := s.getByEmailUserStmt.Close(); err != nil {
		return fmt.Errorf("close get user by email statement: %w", err)
	}
	if err := s.createUserStmt.Close(); err != nil {
		return fmt.Errorf("close create user statement: %w", err)
	}
	if err := s.updateUserStmt.Close(); err != nil {
		return fmt.Errorf("close update user statement: %w", err)
	}
	if err := s.deleteUserStmt.Close(); err != nil {
		return fmt.Errorf("close delete user statement: %w", err)
	}
	if err := s.listTokensStmt.Close(); err != nil {
		return fmt.Errorf("close create user token statement: %w", err)
	}
	if err := s.createTokenStmt.Close(); err != nil {
		return fmt.Errorf("close update user token statement: %w", err)
	}
	if err := s.deleteTokenStmt.Close(); err != nil {
		return fmt.Errorf("close delete user token statement: %w", err)
	}
	if err := s.listRatingsStmt.Close(); err != nil {
		return fmt.Errorf("close list user ratings statement: %w", err)
	}
	if err := s.listUserVehiclesStmt.Close(); err != nil {
		return fmt.Errorf("close list user vehicles statement: %w", err)
	}

	return nil
}

// ListUsers lists all users.
func (s *UserRepository) ListUsers(ctx context.Context) ([]*cargonaut.User, error) {
	users := make([]*cargonaut.User, 0)
	if err := s.listUsersStmt.SelectContext(ctx, &users); err != nil {
		return nil, fmt.Errorf("select users from database: %w", err)
	}
	return users, nil
}

// GetUser returns a user identified by his unique ID.
func (s *UserRepository) GetUser(ctx context.Context, id uuid.UUID) (*cargonaut.User, error) {
	user := new(cargonaut.User)
	if err := s.getUserStmt.GetContext(ctx, user, id); err == sql.ErrNoRows {
		return nil, cargonaut.ErrUserNotFound
	} else if err != nil {
		return nil, fmt.Errorf("get user %q from database: %w", id, err)
	}
	return user, nil
}

// GetUserByEmail returns a user identified by his E-Mail address.
func (s *UserRepository) GetUserByEmail(ctx context.Context, email string) (*cargonaut.User, error) {
	user := new(cargonaut.User)
	if err := s.getByEmailUserStmt.GetContext(ctx, user, email); err == sql.ErrNoRows {
		return nil, cargonaut.ErrUserNotFound
	} else if err != nil {
		return nil, fmt.Errorf("get user %q from database: %w", email, err)
	}
	return user, nil
}

// CreateUser creates a new user.
func (s *UserRepository) CreateUser(ctx context.Context, user *cargonaut.User) error {
	if _, err := s.createUserStmt.ExecContext(ctx, user); isAlreadyExistsError(err) {
		return cargonaut.ErrUserExists
	} else if err != nil {
		return fmt.Errorf("create user in database: %w", err)
	}
	return nil
}

// UpdateUser updates a given user.
func (s *UserRepository) UpdateUser(ctx context.Context, user *cargonaut.User) error {
	if _, err := s.updateUserStmt.ExecContext(ctx, user); isAlreadyExistsError(err) {
		return cargonaut.ErrUserExists
	} else if err != nil {
		return fmt.Errorf("update user %q in database: %w", user.ID, err)
	}
	return nil
}

// DeleteUser deletes a user identified by his unique ID.
func (s *UserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if _, err := s.deleteUserStmt.ExecContext(ctx, id); err != nil {
		return fmt.Errorf("delete user %q from database: %w", id, err)
	}
	return nil
}

// ListTokens lists all authentication tokens for the user identified by his
// unique ID.
func (s *UserRepository) ListTokens(ctx context.Context, userID uuid.UUID) ([]*cargonaut.Token, error) {
	tokens := make([]*cargonaut.Token, 0)
	if err := s.listTokensStmt.SelectContext(ctx, &tokens, userID); err == sql.ErrNoRows {
		return nil, cargonaut.ErrUserNotFound
	} else if err != nil {
		return nil, fmt.Errorf("select user tokens from database: %w", err)
	}
	return tokens, nil
}

// CreateToken creates an authentication token for the user identified by the
// tokens unique user ID.
func (s *UserRepository) CreateToken(ctx context.Context, token *cargonaut.Token) error {
	if _, err := s.createTokenStmt.ExecContext(ctx, token); err != nil {
		if isAlreadyExistsError(err) {
			return cargonaut.ErrTokenExists
		}
		return fmt.Errorf("create user token in database: %w", err)
	}
	return nil
}

// DeleteToken deletes an users authentication token. Token and user are
// identified by their unique IDs.
func (s *UserRepository) DeleteToken(ctx context.Context, userID, tokenID uuid.UUID) error {
	if _, err := s.deleteTokenStmt.ExecContext(ctx, userID, tokenID); err != nil {
		return fmt.Errorf("delete user token from database: %w", err)
	}
	return nil
}

// ListRatings lists all ratings for the user identified by his unique ID.
func (s *UserRepository) ListRatings(ctx context.Context, userID uuid.UUID) ([]*cargonaut.Rating, error) {
	ratings := make([]*cargonaut.Rating, 0)
	if err := s.listRatingsStmt.SelectContext(ctx, &ratings, userID); err == sql.ErrNoRows {
		return nil, cargonaut.ErrUserNotFound
	} else if err != nil {
		return nil, fmt.Errorf("select user ratings from database: %w", err)
	}
	return ratings, nil
}

// ListVehicles lists all vehicles for the user identified by his unique ID.
func (s *UserRepository) ListVehicles(ctx context.Context, userID uuid.UUID) ([]*cargonaut.Vehicle, error) {
	vehicles := make([]*cargonaut.Vehicle, 0)
	if err := s.listUserVehiclesStmt.SelectContext(ctx, &vehicles, userID); err != nil {
		return nil, fmt.Errorf("select vehicles of user %q from database: %w", userID, err)
	}
	return vehicles, nil
}
