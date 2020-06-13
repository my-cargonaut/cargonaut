package main

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"

	"github.com/my-cargonaut/cargonaut/sql"
)

type migrateConfig struct {
	PostgresURL string
}

func migrateCmd(ctx context.Context, args []string, cfg *migrateConfig) error {
	var (
		direction migrate.MigrationDirection
		logPrefix string
	)

	if dir := args[0]; dir == "up" {
		direction = migrate.Up
		logPrefix = "Applied"
	} else if dir == "down" {
		direction = migrate.Down
		logPrefix = "Rewinded"
	} else {
		return fmt.Errorf("invalid argument %q", dir)
	}

	db, err := sqlx.ConnectContext(ctx, "postgres", cfg.PostgresURL)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			logger.Print(err)
		}
	}()

	n, err := sql.Migrate(db, direction)
	if err != nil {
		return fmt.Errorf("migrate database: %w", err)
	} else if n > 0 {
		logger.Printf("%s %d database migrations", logPrefix, n)
	}

	return nil
}
