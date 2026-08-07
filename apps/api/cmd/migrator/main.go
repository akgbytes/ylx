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

	m, err := migrate.New(
		"file://migrations",
		cfg.DatabaseURL,
	)
	if err != nil {
		color.Red("✘ init migrator failed: %v", err)
		return
	}

	defer func() {
		sourceErr, dbErr := m.Close()
		if sourceErr != nil {
			color.Red("✘ migrator close source error: %v", sourceErr)
		}
		if dbErr != nil {
			color.Red("✘ migrator close db error: %v", dbErr)
		}
	}()

	runCommand(m, os.Args)
}

func runCommand(m *migrate.Migrate, args []string) {
	switch args[1] {
	case "up":
		runUp(m)
	case "down":
		runDown(m, args)
	case "version":
		runVersion(m)
	case "goto":
		runGoto(m, args)
	case "force":
		runForce(m, args)
	default:
		color.Yellow("usage: migrator <up|down|version|goto|force>")
	}
}

func runUp(m *migrate.Migrate) {
	err := m.Up()

	switch {
	case err == nil:
		color.Green("✔ migrations applied successfully")
	case errors.Is(err, migrate.ErrNoChange):
		color.Green("✔ all migrations already applied")
	default:
		color.Red("✘ migration failed: %v", err)
	}
}

func runDown(m *migrate.Migrate, args []string) {
	if len(args) < minimumDownCommandArgs {
		color.Yellow("usage: migrator down <steps|all>")
		return
	}

	target := args[2]
	if target == "all" {
		if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			color.Red("✘ rollback failed: %v", err)
			return
		}

		color.Green("✔ all migrations rolled back")
		return
	}

	steps, err := strconv.Atoi(target)
	if err != nil || steps <= 0 {
		color.Red("✘ invalid rollback steps")
		return
	}

	if err := m.Steps(-steps); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		color.Red("✘ rollback failed: %v", err)
		return
	}

	color.Green("✔ rolled back %d migration(s)", steps)
}

func runVersion(m *migrate.Migrate) {
	version, dirty, err := m.Version()
	if err != nil {
		color.Red("✘ failed to get migration version: %v", err)
		return
	}

	color.Cyan("ℹ current version: %d (dirty: %t)", version, dirty)
}

func runGoto(m *migrate.Migrate, args []string) {
	gotoCmd := flag.NewFlagSet("goto", flag.ContinueOnError)
	version := gotoCmd.Uint("version", 0, "target migration version")

	if err := gotoCmd.Parse(args[2:]); err != nil {
		color.Red("✘ invalid goto arguments: %v", err)
		return
	}
	if *version == 0 {
		color.Yellow("usage: migrator goto --version <version>")
		return
	}
	if err := m.Migrate(*version); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		color.Red("✘ migration failed: %v", err)
		return
	}

	color.Green("✔ migrated to version %d", *version)
}

func runForce(m *migrate.Migrate, args []string) {
	forceCmd := flag.NewFlagSet("force", flag.ContinueOnError)
	version := forceCmd.Uint("version", 0, "force migration version")

	if err := forceCmd.Parse(args[2:]); err != nil {
		color.Red("✘ invalid force arguments: %v", err)
		return
	}
	if *version == 0 {
		color.Yellow("usage: migrator force --version <version>")
		return
	}
	if err := m.Force(int(*version)); err != nil {
		color.Red("✘ force failed: %v", err)
		return
	}

	color.Green("✔ forced migration version to %d", *version)
}
