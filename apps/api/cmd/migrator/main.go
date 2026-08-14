package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	"github.com/rs/zerolog"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/akgbytes/ylx/internal/config"
	"github.com/akgbytes/ylx/internal/logger"
)

const (
	minCommandArguments = 2
	usage               = "usage: migrator <up|down|version|goto|force>"
)

func main() {
	log := logger.BootstrapLogger()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("bootstrap migrator")
	}

	if len(os.Args) < minCommandArguments {
		log.Fatal().Msg(usage)
	}

	migrationsPath, err := filepath.Abs("migrations")
	if err != nil {
		log.Fatal().Err(err).Msg("resolve migrations directory")
	}

	migrator, err := migrate.New("file://"+migrationsPath, cfg.Database.URL)
	if err != nil {
		log.Fatal().Err(err).Msg("initialize migrator")
	}
	if err := runMigrations(migrator, os.Args[1:], log); err != nil {
		log.Fatal().Err(err).Msg("migration failed")
	}
}

func runMigrations(migrator *migrate.Migrate, args []string, log zerolog.Logger) error {
	defer func() {
		sourceErr, databaseErr := migrator.Close()
		if err := errors.Join(sourceErr, databaseErr); err != nil {
			log.Error().Err(err).Msg("close migrator")
		}
	}()

	if err := run(migrator, args); err != nil {
		return err
	}

	if args[0] != "version" {
		log.Info().Msg("migration completed")
	}

	return nil
}

func run(migrator *migrate.Migrate, args []string) error {
	switch args[0] {
	case "up":
		return noChangeIsNil(migrator.Up())
	case "down":
		return down(migrator, args[1:])
	case "version":
		version, dirty, err := migrator.Version()
		if errors.Is(err, migrate.ErrNilVersion) {
			return writeVersion("version: none\n")
		}
		if err != nil {
			return fmt.Errorf("get version: %w", err)
		}
		return writeVersion("version: %d (dirty: %t)\n", version, dirty)
	case "goto":
		return goTo(migrator, args[1:])
	case "force":
		return force(migrator, args[1:])
	default:
		return errors.New(usage)
	}
}

func writeVersion(format string, args ...any) error {
	if _, err := fmt.Fprintf(os.Stdout, format, args...); err != nil {
		return fmt.Errorf("write version: %w", err)
	}
	return nil
}

func down(migrator *migrate.Migrate, args []string) error {
	if len(args) != 1 {
		return errors.New("usage: migrator down <steps|all>")
	}
	if args[0] == "all" {
		return noChangeIsNil(migrator.Down())
	}

	steps, err := strconv.Atoi(args[0])
	if err != nil || steps <= 0 {
		return errors.New("rollback steps must be a positive integer")
	}
	return noChangeIsNil(migrator.Steps(-steps))
}

func goTo(migrator *migrate.Migrate, args []string) error {
	flags := flag.NewFlagSet("goto", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	version := flags.Uint("version", 0, "target migration version")
	if err := flags.Parse(args); err != nil {
		return fmt.Errorf("parse goto arguments: %w", err)
	}
	if *version == 0 {
		return errors.New("usage: migrator goto --version <version>")
	}
	return noChangeIsNil(migrator.Migrate(*version))
}

func force(migrator *migrate.Migrate, args []string) error {
	flags := flag.NewFlagSet("force", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	version := flags.Uint("version", 0, "force migration version")
	if err := flags.Parse(args); err != nil {
		return fmt.Errorf("parse force arguments: %w", err)
	}
	if *version == 0 {
		return errors.New("usage: migrator force --version <version>")
	}
	return migrator.Force(int(*version))
}

func noChangeIsNil(err error) error {
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}
	return err
}
