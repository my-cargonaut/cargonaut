package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"

	"github.com/my-cargonaut/cargonaut/http"
	"github.com/my-cargonaut/cargonaut/sql"
)

type serveConfig struct {
	Automigrate   bool
	ListenAddress string
	PostgresURL   string
	SecureCookies bool
}

func serveCmd(ctx context.Context, _ []string, cfg *serveConfig) error {
	db, err := sqlx.ConnectContext(ctx, "postgres", cfg.PostgresURL)
	if err != nil {
		return fmt.Errorf("create database: %w", err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			logger.Print(err)
		}
	}()
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Minute * 5)

	if cfg.Automigrate {
		var n int
		if n, err = sql.Migrate(db, migrate.Up); err != nil {
			return fmt.Errorf("migrate database: %w", err)
		} else if n > 0 {
			logger.Printf("Applied %d database migrations", n)
		}
	}

	userService, err := sql.NewUserService(ctx, db)
	if err != nil {
		return fmt.Errorf("create truck service: %w", err)
	}
	defer func() {
		if err = userService.Close(); err != nil {
			logger.Print(err)
		}
	}()

	h, err := http.NewHandler(logger)
	if err != nil {
		return fmt.Errorf("create http handler: %w", err)
	}
	h.UserService = userService

	srv, err := http.NewServer(logger, cfg.ListenAddress, h)
	if err != nil {
		return fmt.Errorf("create http server: %w", err)
	}
	srv.Run()

	logger.Printf("Listening on %s", srv.ListenAddr().String())

	select {
	case <-ctx.Done():
		logger.Print("Gracefully shutting down server")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Print(err)
		}
	case err := <-srv.ListenError():
		return fmt.Errorf("start http server: %w", err)
	}

	return nil
}
