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

var _ cargonaut.VehicleRepository = (*VehicleRepository)(nil)

const (
	listVehiclesSQL  = "SELECT id, user_id, brand, model, passengers, loading_area_length, loading_area_width, created_at, updated_at FROM vehicle ORDER BY updated_at DESC"
	getVehicleSQL    = "SELECT id, user_id, brand, model, passengers, loading_area_length, loading_area_width, created_at, updated_at FROM vehicle WHERE id = $1 LIMIT 1"
	createVehicleSQL = "INSERT INTO vehicle (user_id, brand, model, passengers, loading_area_length, loading_area_width) VALUES (:user_id, :brand, :model, :passengers, :loading_area_length, :loading_area_width)"
	updateVehicleSQL = "UPDATE vehicle SET brand = :brand, model = :model, passengers = :passengers, loading_area_length = :loading_area_length, loading_area_width = :loading_area_width, updated_at = :updated_at WHERE id = :id"
	deleteVehicleSQL = "DELETE FROM vehicle WHERE id = $1"
)

// VehicleRepository provides access to the vehicle resource backed by a Postgres
// SQL database.
type VehicleRepository struct {
	db *sqlx.DB

	listStmt   *sqlx.Stmt
	getStmt    *sqlx.Stmt
	createStmt *sqlx.NamedStmt
	updateStmt *sqlx.NamedStmt
	deleteStmt *sqlx.Stmt
}

// NewVehicleRepository returns a new VehicleRepository based on top of the
// provided database connection.
func NewVehicleRepository(ctx context.Context, db *sqlx.DB) (*VehicleRepository, error) {
	s := &VehicleRepository{db: db}

	var err error
	if s.listStmt, err = db.PreparexContext(ctx, listVehiclesSQL); err != nil {
		return nil, fmt.Errorf("prepare list vehicles statement: %w", err)
	}
	if s.getStmt, err = db.PreparexContext(ctx, getVehicleSQL); err != nil {
		return nil, fmt.Errorf("prepare get vehicle statement: %w", err)
	}
	if s.createStmt, err = db.PrepareNamedContext(ctx, createVehicleSQL); err != nil {
		return nil, fmt.Errorf("prepare create vehicle statement: %w", err)
	}
	if s.updateStmt, err = db.PrepareNamedContext(ctx, updateVehicleSQL); err != nil {
		return nil, fmt.Errorf("prepare update vehicle statement: %w", err)
	}
	if s.deleteStmt, err = db.PreparexContext(ctx, deleteVehicleSQL); err != nil {
		return nil, fmt.Errorf("prepare delete vehicle statement: %w", err)
	}

	return s, nil
}

// Close all prepared statements.
func (s *VehicleRepository) Close() error {
	if err := s.listStmt.Close(); err != nil {
		return fmt.Errorf("close list vehicle statement: %w", err)
	}
	if err := s.getStmt.Close(); err != nil {
		return fmt.Errorf("close get vehicle statement: %w", err)
	}
	if err := s.createStmt.Close(); err != nil {
		return fmt.Errorf("close create vehicle statement: %w", err)
	}
	if err := s.updateStmt.Close(); err != nil {
		return fmt.Errorf("close update vehicle statement: %w", err)
	}
	if err := s.deleteStmt.Close(); err != nil {
		return fmt.Errorf("close delete vehicle statement: %w", err)
	}

	return nil
}

// ListVehicles lists all vehicles.
func (s *VehicleRepository) ListVehicles(ctx context.Context) ([]*cargonaut.Vehicle, error) {
	vehicles := make([]*cargonaut.Vehicle, 0)
	if err := s.listStmt.SelectContext(ctx, &vehicles); err != nil {
		return nil, fmt.Errorf("select vehicles from database: %w", err)
	}
	return vehicles, nil
}

// GetVehicle returns a vehicle identified by his unique ID.
func (s *VehicleRepository) GetVehicle(ctx context.Context, id uuid.UUID) (*cargonaut.Vehicle, error) {
	vehicle := new(cargonaut.Vehicle)
	if err := s.getStmt.GetContext(ctx, vehicle, id); err == sql.ErrNoRows {
		return nil, cargonaut.ErrVehicleNotFound
	} else if err != nil {
		return nil, fmt.Errorf("get vehicle %q from database: %w", id, err)
	}
	return vehicle, nil
}

// CreateVehicle creates a new vehicle.
func (s *VehicleRepository) CreateVehicle(ctx context.Context, vehicle *cargonaut.Vehicle) error {
	if _, err := s.createStmt.ExecContext(ctx, vehicle); isAlreadyExistsError(err) {
		return cargonaut.ErrVehicleExists
	} else if err != nil {
		return fmt.Errorf("create vehicle in database: %w", err)
	}
	return nil
}

// UpdateVehicle updates a given vehicle.
func (s *VehicleRepository) UpdateVehicle(ctx context.Context, vehicle *cargonaut.Vehicle) error {
	if _, err := s.updateStmt.ExecContext(ctx, vehicle); isAlreadyExistsError(err) {
		return cargonaut.ErrVehicleExists
	} else if err != nil {
		return fmt.Errorf("update vehicle %q in database: %w", vehicle.ID, err)
	}
	return nil
}

// DeleteVehicle deletes a vehicle identified by his unique ID.
func (s *VehicleRepository) DeleteVehicle(ctx context.Context, id uuid.UUID) error {
	if _, err := s.deleteStmt.ExecContext(ctx, id); err != nil {
		return fmt.Errorf("delete vehicle %q from database: %w", id, err)
	}
	return nil
}
