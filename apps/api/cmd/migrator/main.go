package main

import (
	"errors"
	"flag"
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/akgbytes/ylx/internal/config"
)

const (
	minimumCommandArgs     = 2
	minimumDownCommandArgs = 3
)

func main() {
	cfg := config.MustLoad()

	if len(os.Args) < minimumCommandArgs {
		color.Yellow("usage: migrator <up|down|version|goto|force>")
		return
	}

	migrator, err := migrate.New(
		"file://migrations",
		cfg.DatabaseURL,
	)
	if err != nil {
		color.Red("✘ failed to initialize migrator: %v", err)
		return
	}

	defer func() {
		sourceErr, databaseErr := migrator.Close()
		if sourceErr != nil {
			color.Red("✘ failed to close migration source: %v", sourceErr)
		}
		if databaseErr != nil {
			color.Red("✘ failed to close migration database: %v", databaseErr)
		}
	}()

	dispatchCommand(migrator, os.Args)
}

func dispatchCommand(migrator *migrate.Migrate, args []string) {
	switch args[1] {
	case "up":
		applyMigrations(migrator)
	case "down":
		rollbackMigrations(migrator, args)
	case "version":
		showVersion(migrator)
	case "goto":
		migrateToVersion(migrator, args)
	case "force":
		forceVersion(migrator, args)
	default:
		color.Yellow("usage: migrator <up|down|version|goto|force>")
	}
}

func applyMigrations(migrator *migrate.Migrate) {
	err := migrator.Up()

	switch {
	case err == nil:
		color.Green("✔ migrations applied successfully")
	case errors.Is(err, migrate.ErrNoChange):
		color.Green("✔ all migrations already applied")
	default:
		color.Red("✘ migration failed: %v", err)
	}
}

func rollbackMigrations(migrator *migrate.Migrate, args []string) {
	if len(args) < minimumDownCommandArgs {
		color.Yellow("usage: migrator down <steps|all>")
		return
	}

	stepsOrAll := args[2]
	if stepsOrAll == "all" {
		if err := migrator.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			color.Red("✘ rollback failed: %v", err)
			return
		}

		color.Green("✔ all migrations rolled back")
		return
	}

	steps, err := strconv.Atoi(stepsOrAll)
	if err != nil || steps <= 0 {
		color.Red("✘ rollback steps must be a positive integer")
		return
	}

	if err := migrator.Steps(-steps); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		color.Red("✘ rollback failed: %v", err)
		return
	}

	color.Green("✔ rolled back %d migration(s)", steps)
}

func showVersion(migrator *migrate.Migrate) {
	version, dirty, err := migrator.Version()
	if err != nil {
		color.Red("✘ failed to get migration version: %v", err)
		return
	}

	color.Cyan("ℹ current version: %d (dirty: %t)", version, dirty)
}

func migrateToVersion(migrator *migrate.Migrate, args []string) {
	gotoFlags := flag.NewFlagSet("goto", flag.ContinueOnError)
	version := gotoFlags.Uint("version", 0, "target migration version")

	if err := gotoFlags.Parse(args[2:]); err != nil {
		color.Red("✘ invalid goto arguments: %v", err)
		return
	}
	if *version == 0 {
		color.Yellow("usage: migrator goto --version <version>")
		return
	}
	if err := migrator.Migrate(*version); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		color.Red("✘ migration failed: %v", err)
		return
	}

	color.Green("✔ migrated to version %d", *version)
}

func forceVersion(migrator *migrate.Migrate, args []string) {
	forceFlags := flag.NewFlagSet("force", flag.ContinueOnError)
	version := forceFlags.Uint("version", 0, "force migration version")

	if err := forceFlags.Parse(args[2:]); err != nil {
		color.Red("✘ invalid force arguments: %v", err)
		return
	}
	if *version == 0 {
		color.Yellow("usage: migrator force --version <version>")
		return
	}
	if err := migrator.Force(int(*version)); err != nil {
		color.Red("✘ force failed: %v", err)
		return
	}

	color.Green("✔ forced migration version to %d", *version)
}
