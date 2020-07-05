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

var _ cargonaut.TripRepository = (*TripRepository)(nil)

const (
	listTripsSQL    = "SELECT id, user_id, vehicle_id, rider_id, start, destination, price, depature, arrival, created_at, updated_at FROM trip ORDER BY updated_at DESC"
	getTripSQL      = "SELECT id, user_id, vehicle_id, rider_id, start, destination, price, depature, arrival, created_at, updated_at FROM trip WHERE id = $1 LIMIT 1"
	createTripSQL   = "INSERT INTO trip (user_id, vehicle_id, rider_id, start, destination, price, depature, arrival) VALUES (:user_id, :vehicle_id, :rider_id, :start, :destination, :price, :depature, :arrival)"
	updateTripSQL   = "UPDATE trip SET vehicle_id = :vehicle_id, rider_id = :rider_id, start = :start, destination = :destination, price = :price, depature = :depature, arrival = :arrival, updated_at = :updated_at WHERE id = :id"
	deleteTripSQL   = "DELETE FROM trip WHERE id = $1"
	getRatingSQL    = "SELECT id, user_id, author_id, trip_id, value, comment, created_at FROM rating WHERE trip_id = $1 LIMIT 1"
	createRatingSQL = "INSERT INTO rating (user_id, author_id, trip_id, comment, value) VALUES (:user_id, :author_id, :trip_id, :comment, :value)"
)

// TripRepository provides access to the trip resource backed by a Postgres SQL
// database.
type TripRepository struct {
	db *sqlx.DB

	listTripsStmt    *sqlx.Stmt
	getStmt          *sqlx.Stmt
	createStmt       *sqlx.NamedStmt
	updateStmt       *sqlx.NamedStmt
	deleteStmt       *sqlx.Stmt
	getRatingStmt    *sqlx.Stmt
	createRatingStmt *sqlx.NamedStmt
}

// NewTripRepository returns a new TripRepository based on top of the provided
// database connection.
func NewTripRepository(ctx context.Context, db *sqlx.DB) (*TripRepository, error) {
	s := &TripRepository{db: db}

	var err error
	if s.listTripsStmt, err = db.PreparexContext(ctx, listTripsSQL); err != nil {
		return nil, fmt.Errorf("prepare list trips statement: %w", err)
	}
	if s.getStmt, err = db.PreparexContext(ctx, getTripSQL); err != nil {
		return nil, fmt.Errorf("prepare get trip statement: %w", err)
	}
	if s.createStmt, err = db.PrepareNamedContext(ctx, createTripSQL); err != nil {
		return nil, fmt.Errorf("prepare create trip statement: %w", err)
	}
	if s.updateStmt, err = db.PrepareNamedContext(ctx, updateTripSQL); err != nil {
		return nil, fmt.Errorf("prepare update trip statement: %w", err)
	}
	if s.deleteStmt, err = db.PreparexContext(ctx, deleteTripSQL); err != nil {
		return nil, fmt.Errorf("prepare delete trip statement: %w", err)
	}
	if s.getRatingStmt, err = db.PreparexContext(ctx, getRatingSQL); err != nil {
		return nil, fmt.Errorf("prepare get trip rating statement: %w", err)
	}
	if s.createRatingStmt, err = db.PrepareNamedContext(ctx, createRatingSQL); err != nil {
		return nil, fmt.Errorf("prepare create trip rating statement: %w", err)
	}

	return s, nil
}

// Close all prepared statements.
func (s *TripRepository) Close() error {
	if err := s.listTripsStmt.Close(); err != nil {
		return fmt.Errorf("close list trip statement: %w", err)
	}
	if err := s.getStmt.Close(); err != nil {
		return fmt.Errorf("close get trip statement: %w", err)
	}
	if err := s.createStmt.Close(); err != nil {
		return fmt.Errorf("close create trip statement: %w", err)
	}
	if err := s.updateStmt.Close(); err != nil {
		return fmt.Errorf("close update trip statement: %w", err)
	}
	if err := s.deleteStmt.Close(); err != nil {
		return fmt.Errorf("close delete trip statement: %w", err)
	}
	if err := s.getRatingStmt.Close(); err != nil {
		return fmt.Errorf("close get trip rating statement: %w", err)
	}
	if err := s.createRatingStmt.Close(); err != nil {
		return fmt.Errorf("close create trip rating statement: %w", err)
	}

	return nil
}

// ListTrips lists all trips.
func (s *TripRepository) ListTrips(ctx context.Context) ([]*cargonaut.Trip, error) {
	trips := make([]*cargonaut.Trip, 0)
	if err := s.listTripsStmt.SelectContext(ctx, &trips); err != nil {
		return nil, fmt.Errorf("select trips from database: %w", err)
	}
	return trips, nil
}

// GetTrip returns a trip identified by his unique ID.
func (s *TripRepository) GetTrip(ctx context.Context, id uuid.UUID) (*cargonaut.Trip, error) {
	trip := new(cargonaut.Trip)
	if err := s.getStmt.GetContext(ctx, trip, id); err == sql.ErrNoRows {
		return nil, cargonaut.ErrTripNotFound
	} else if err != nil {
		return nil, fmt.Errorf("get trip %q from database: %w", id, err)
	}
	return trip, nil
}

// CreateTrip creates a new trip.
func (s *TripRepository) CreateTrip(ctx context.Context, trip *cargonaut.Trip) error {
	if _, err := s.createStmt.ExecContext(ctx, trip); isAlreadyExistsError(err) {
		return cargonaut.ErrTripExists
	} else if err != nil {
		return fmt.Errorf("create trip in database: %w", err)
	}
	return nil
}

// UpdateTrip updates a given trip.
func (s *TripRepository) UpdateTrip(ctx context.Context, trip *cargonaut.Trip) error {
	if _, err := s.updateStmt.ExecContext(ctx, trip); isAlreadyExistsError(err) {
		return cargonaut.ErrTripExists
	} else if err != nil {
		return fmt.Errorf("update trip %q in database: %w", trip.ID, err)
	}
	return nil
}

// DeleteTrip deletes a trip identified by his unique ID.
func (s *TripRepository) DeleteTrip(ctx context.Context, id uuid.UUID) error {
	if _, err := s.deleteStmt.ExecContext(ctx, id); err != nil {
		return fmt.Errorf("delete trip %q from database: %w", id, err)
	}
	return nil
}

// GetRating gets a rating for a trip.
func (s *TripRepository) GetRating(ctx context.Context, id uuid.UUID) (*cargonaut.Rating, error) {
	rating := new(cargonaut.Rating)
	if err := s.getRatingStmt.GetContext(ctx, rating, id); err == sql.ErrNoRows {
		return nil, cargonaut.ErrRatingNotFound
	} else if err != nil {
		return nil, fmt.Errorf("get rating for trip %q from database: %w", id, err)
	}
	return rating, nil
}

// CreateRating creates a new rating fro a trip.
func (s *TripRepository) CreateRating(ctx context.Context, rating *cargonaut.Rating) error {
	if _, err := s.createRatingStmt.ExecContext(ctx, rating); err != nil {
		if isAlreadyExistsError(err) {
			return cargonaut.ErrRatingExists
		}
		return fmt.Errorf("create rating for trip %q in database: %w", rating.TripID, err)
	}
	return nil
}
