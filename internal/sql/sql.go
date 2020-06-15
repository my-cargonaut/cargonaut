package sql

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	migrate "github.com/rubenv/sql-migrate"
)

// Migrate database schema into the given direction.
func Migrate(db *sqlx.DB, direction migrate.MigrationDirection) (int, error) {
	migrations, err := fs.NewWithNamespace("migrations")
	if err != nil {
		return 0, fmt.Errorf("create migrations file system: %w", err)
	}

	n, err := migrate.Exec(db.DB, "postgres", &migrate.HttpFileSystemMigrationSource{
		FileSystem: migrations,
	}, direction)
	if err != nil {
		return 0, fmt.Errorf("apply database migrations: %w", err)
	}
	return n, nil
}

// isAlreadyExistsError returns true if the supplied error indicates that a
// resource with the same constraints already exists.
func isAlreadyExistsError(err error) bool {
	if err == nil {
		return false
	}
	psqlErr, ok := err.(*pq.Error)
	return ok && psqlErr.Code == "23505"
}
