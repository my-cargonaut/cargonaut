package main

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"

	"github.com/my-cargonaut/cargonaut/internal/handler"
	"github.com/my-cargonaut/cargonaut/internal/redis"
	"github.com/my-cargonaut/cargonaut/internal/sql"
	"github.com/my-cargonaut/cargonaut/pkg/http"
)

type serveConfig struct {
	Automigrate   bool
	ListenAddress string
	PostgresURL   string
	RedisURL      string
	Secret        string
}

func serveCmd(ctx context.Context, _ []string, cfg *serveConfig) error {
	// Decode the hex encoded secret.
	secret, err := hex.DecodeString(cfg.Secret)
	if err != nil {
		return fmt.Errorf("decode secret: %w", err)
	} else if len(secret) != 32 {
		return errors.New("secret must be 32 bytes long")
	}

	// Connect to PostgreSQL database.
	db, err := sqlx.ConnectContext(ctx, "postgres", cfg.PostgresURL)
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			logger.Printf("close database connection: %s", err)
		}
	}()
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Minute * 5)

	// Connect to the Redis cache.
	cache, err := redigo.DialURL(cfg.RedisURL,
		redigo.DialConnectTimeout(time.Second*5),
		redigo.DialReadTimeout(time.Second*5),
		redigo.DialWriteTimeout(time.Second*5),
		redigo.DialKeepAlive(time.Minute*5),
	)
	if err != nil {
		return fmt.Errorf("connect to cache: %w", err)
	}
	defer func() {
		if err = cache.Close(); err != nil {
			logger.Printf("close cache connection: %s", err)
		}
	}()

	// Check connection to Redis server.
	if _, err = cache.Do("PING"); err != nil {
		logger.Printf("check cache connection: %s", err)
	}

	// Run database migrations, if auto migrate is set.
	if cfg.Automigrate {
		var n int
		if n, err = sql.Migrate(db, migrate.Up); err != nil {
			return fmt.Errorf("migrate database: %w", err)
		} else if n > 0 {
			logger.Printf("Applied %d database migrations", n)
		}
	}

	userRepository, err := sql.NewUserRepository(ctx, db)
	if err != nil {
		return fmt.Errorf("create user repository: %w", err)
	}
	defer func() {
		if err = userRepository.Close(); err != nil {
			logger.Printf("close user repository: %s", err)
		}
	}()

	tokenBlacklist := redis.NewTokenBlacklist(cache)

	// Create http handlers.
	h, err := handler.NewHandler(logger, secret)
	if err != nil {
		return fmt.Errorf("create http handler: %w", err)
	}
	h.UserRepository = userRepository
	h.TokenBlacklist = tokenBlacklist

	// Run http server.
	srv, err := http.NewServer(logger, cfg.ListenAddress, h)
	if err != nil {
		return fmt.Errorf("create http server: %w", err)
	}
	srv.Run()

	logger.Printf("Listening on %s", srv.ListenAddr().String())

	// Wait for cancellation or http server error.
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
