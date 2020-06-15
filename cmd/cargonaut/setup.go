package main

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"

	"github.com/my-cargonaut/cargonaut"
	"github.com/my-cargonaut/cargonaut/internal/sql"
	"github.com/my-cargonaut/cargonaut/pkg/password"
	"github.com/my-cargonaut/cargonaut/pkg/prompt"
)

type setupConfig struct {
	PostgresURL string
	Secret      string
}

type setupUser struct {
	Email       string `prompt:"required,name=E-Mail,validate=email"`
	Password    string `prompt:"required,mask=*,validate=password"`
	DisplayName string `prompt:"required,name=Display Name"`
}

func setupCmd(ctx context.Context, _ []string, cfg *setupConfig) error {
	// Decode the hex encoded secret.
	secret, err := hex.DecodeString(cfg.Secret)
	if err != nil {
		return fmt.Errorf("decode secret: %w", err)
	} else if len(secret) != 32 {
		return errors.New("secret must be 32 bytes long")
	}

	// Connect to database.
	db, err := sqlx.ConnectContext(ctx, "postgres", cfg.PostgresURL)
	if err != nil {
		return fmt.Errorf("create database: %w", err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			logger.Print(err)
		}
	}()

	// Run database migrations.
	var n int
	if n, err = sql.Migrate(db, migrate.Up); err != nil {
		return fmt.Errorf("migrate database: %w", err)
	} else if n > 0 {
		logger.Printf("Applied %d database migrations", n)
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

	// Prompt for user input.
	var userSetup setupUser
	if err = prompt.NewWithValidators(defaultPromptValidators).Run(&userSetup); err != nil {
		return fmt.Errorf("prompt for user input: %w", err)
	}

	user := &cargonaut.User{
		Email:       userSetup.Email,
		Password:    userSetup.Password,
		DisplayName: userSetup.DisplayName,
	}

	// Hash supplied password and create the user resource.
	if user.Password, err = password.Generate(secret, user.Password, password.DefaultCost); err != nil {
		return fmt.Errorf("hash user password: %w", err)
	} else if err = userService.CreateUser(ctx, user); err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}
