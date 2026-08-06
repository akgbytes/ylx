package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/akgbytes/ylx/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg := config.MustLoad()

	if len(os.Args) < 2 {
		fmt.Println("Usage: migrator <up|down|version|goto|force>")
		return
	}

	m, err := migrate.New(
		"file://migrations",
		cfg.DatabaseURL,
	)
	if err != nil {
		fmt.Println("Init migrator failed:", err)
		return
	}

	defer m.Close()

	switch os.Args[1] {
	case "up":
		err := m.Up()

		switch err {
		case nil:
			fmt.Println("Migrations applied successfully")
		case migrate.ErrNoChange:
			fmt.Println("All migrations already applied")
		default:
			fmt.Println("Migration failed:", err)
		}

	case "down":
		if len(os.Args) < 3 {
			fmt.Println("Usage: migrator down <steps|all>")
			return
		}

		target := os.Args[2]

		if target == "all" {
			if err := m.Down(); err != nil && err != migrate.ErrNoChange {
				fmt.Println("Rollback failed:", err)
				return
			}

			fmt.Println("All migrations rolled back")
			return
		}

		steps, err := strconv.Atoi(target)
		if err != nil || steps <= 0 {
			fmt.Println("Invalid rollback steps")
			return
		}

		if err := m.Steps(-steps); err != nil && err != migrate.ErrNoChange {
			fmt.Println("Rollback failed:", err)
			return
		}

		fmt.Printf("Rolled back %d migration(s)\n", steps)

	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			fmt.Println("Failed to get migration version:", err)
			return
		}

		fmt.Printf("Current version: %d (dirty: %t)\n", version, dirty)

	case "goto":
		gotoCmd := flag.NewFlagSet("goto", flag.ExitOnError)

		version := gotoCmd.Uint(
			"version",
			0,
			"target migration version",
		)

		gotoCmd.Parse(os.Args[2:])

		if *version == 0 {
			fmt.Println("Usage: migrator goto --version <version>")
			return
		}

		if err := m.Migrate(*version); err != nil && err != migrate.ErrNoChange {
			fmt.Println("Migration failed:", err)
			return
		}

		fmt.Printf("Migrated to version %d\n", *version)

	case "force":
		forceCmd := flag.NewFlagSet("force", flag.ExitOnError)

		version := forceCmd.Uint(
			"version",
			0,
			"force migration version",
		)

		forceCmd.Parse(os.Args[2:])

		if *version == 0 {
			fmt.Println("Usage: migrator force --version <version>")
			return
		}

		if err := m.Force(int(*version)); err != nil {
			fmt.Println("Force failed:", err)
			return
		}

		fmt.Printf("Forced migration version to %d\n", *version)

	default:
		fmt.Println("Usage: migrator <up|down|version|goto|force>")
	}
}
