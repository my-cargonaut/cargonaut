package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"github.com/peterbourgon/ff/v2"
	"github.com/peterbourgon/ff/v2/ffcli"

	"github.com/my-cargonaut/cargonaut/pkg/prompt"
	"github.com/my-cargonaut/cargonaut/pkg/version"
)

var logger = log.New(os.Stdout, "", 0)

var defaultPromptValidators = map[string]prompt.ValidatorFunc{
	"email":    validateEmail,
	"password": validatePassword,
}

func main() {
	// Set up signal handling.
	term := make(chan os.Signal, 1)
	defer close(term)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(term)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		select {
		case <-ctx.Done():
		case <-term:
			logger.Print("Received interrupt")
			cancel()
		}
	}()

	var (
		migrateCfg migrateConfig
		serveCfg   serveConfig
		setupCfg   setupConfig
	)

	migrate := &ffcli.Command{
		Name:       "migrate",
		ShortUsage: "migrate [flags] (up|down)",
		ShortHelp:  "Run the database migrations",
		FlagSet:    flag.NewFlagSet("migrate", flag.ExitOnError),
		Options: []ff.Option{
			ff.WithEnvVarPrefix("CARGONAUT"),
		},
		Exec: func(ctx context.Context, args []string) error {
			if len(args) != 1 {
				return flag.ErrHelp
			}
			return migrateCmd(ctx, args, &migrateCfg)
		},
	}
	serve := &ffcli.Command{
		Name:       "serve",
		ShortUsage: "serve [flags]",
		ShortHelp:  "Serve the HTTP server",
		FlagSet:    flag.NewFlagSet("serve", flag.ExitOnError),
		Options: []ff.Option{
			ff.WithEnvVarPrefix("CARGONAUT"),
		},
		Exec: func(ctx context.Context, args []string) error {
			return serveCmd(ctx, args, &serveCfg)
		},
	}
	setup := &ffcli.Command{
		Name:       "setup",
		ShortUsage: "setup [flags]",
		ShortHelp:  "Setup the application and create a new user",
		FlagSet:    flag.NewFlagSet("setup", flag.ExitOnError),
		Options: []ff.Option{
			ff.WithEnvVarPrefix("CARGONAUT"),
		},
		Exec: func(ctx context.Context, args []string) error {
			return setupCmd(ctx, args, &setupCfg)
		},
	}
	version := &ffcli.Command{
		Name:       "version",
		ShortUsage: "version",
		ShortHelp:  "Print the cargonaut version and build details",
		Exec: func(ctx context.Context, args []string) error {
			fmt.Println(version.Print("cargonaut"))
			return nil
		},
	}
	root := &ffcli.Command{
		ShortUsage: "cargonaut <subcommand>",
		ShortHelp:  "Simple, scalable, fast fuel tank management",
		LongHelp: `cargonaut is a simple, scalable and fast fuel
tank management service.

> Documentation & Support: https://github.com/my-cargonaut/cargonaut
> Source & Copyright Information: https://github.com/my-cargonaut/cargonaut`,
		Subcommands: []*ffcli.Command{migrate, serve, setup, version},
		Exec: func(ctx context.Context, args []string) error {
			return serve.ParseAndRun(ctx, args)
		},
	}

	migrate.FlagSet.StringVar(&migrateCfg.PostgresURL, "postgres-url", "", "URL of the Postgres instance")
	serve.FlagSet.BoolVar(&serveCfg.Automigrate, "automigrate", false, "automatically run database migrations")
	serve.FlagSet.StringVar(&serveCfg.ListenAddress, "listen-address", "", "listen address")
	serve.FlagSet.StringVar(&serveCfg.PostgresURL, "postgres-url", "", "URL of the Postgres instance")
	serve.FlagSet.StringVar(&serveCfg.RedisURL, "redis-url", "", "URL of the Redis instance")
	serve.FlagSet.StringVar(&serveCfg.Secret, "secret", "", "Hex encoded 32 byte secret key for AES-128 GCM and HS256")
	setup.FlagSet.StringVar(&setupCfg.PostgresURL, "postgres-url", "", "URL of the Postgres instance")
	setup.FlagSet.StringVar(&setupCfg.Secret, "secret", "", "Hex encoded 32 byte secret key for AES-128 GCM and HS256")

	if err := root.ParseAndRun(ctx, os.Args[1:]); err != nil && err != flag.ErrHelp {
		logger.Print(err)
		os.Exit(1)
	}
}

var validEmailRegexp = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

func validateEmail(input string) error {
	if !validEmailRegexp.MatchString(input) {
		return errors.New("not a valid email address")
	}
	return nil
}

func validatePassword(input string) error {
	if l := len(input); l < 8 {
		return fmt.Errorf("must be at least 8 characters but is only %d", l)
	}
	return nil
}
