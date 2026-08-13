package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/fatih/color"

	"github.com/akgbytes/ylx/internal/config"
	"github.com/akgbytes/ylx/internal/database"
)

const (
	defaultSeedDir         = "seeds"
	seedTransactionTimeout = 30 * time.Second
)

func main() {
	if err := run(); err != nil {
		color.Red("✘ %v", err)
		os.Exit(1)
	}
}

func run() error {
	seedDir := flag.String("dir", defaultSeedDir, "directory containing seed files")
	seedFile := flag.String("file", "", "seed file to execute")
	flag.Parse()

	cfg := config.MustLoad()

	db, err := database.Connect(context.Background(), *cfg)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			color.Red("✘ failed to close database: %v", err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), seedTransactionTimeout)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	// Safe even after Commit(). Rollback() will simply return sql.ErrTxDone.
	defer func() {
		_ = tx.Rollback()
	}()

	if err := executeSeeds(ctx, tx, *seedDir, *seedFile); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	color.Green("✔ database seeded")

	return nil
}

func executeSeeds(ctx context.Context, tx *sql.Tx, seedDir, seedFile string) error {
	if seedFile == "" {
		return executeSeedDirectory(ctx, tx, seedDir)
	}

	if filepath.Ext(seedFile) != ".sql" {
		return errors.New("seed file must have a .sql extension")
	}

	path := seedFile

	// Resolve relative paths against the configured seed directory.
	// Absolute paths are used as-is.
	if !filepath.IsAbs(path) {
		path = filepath.Join(seedDir, path)
	}

	return executeSeedFile(ctx, tx, path)
}

func executeSeedDirectory(ctx context.Context, tx *sql.Tx, seedDir string) error {
	entries, err := os.ReadDir(seedDir)
	if err != nil {
		return fmt.Errorf("read seed directory: %w", err)
	}

	var files []string

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".sql" {
			continue
		}

		files = append(files, filepath.Join(seedDir, entry.Name()))
	}

	if len(files) == 0 {
		return fmt.Errorf("no .sql seed files found in %q", seedDir)
	}

	// Execute in alphabetical order (001_, 002_, 003_, ...)
	sort.Strings(files)

	for _, seedFile := range files {
		color.Cyan("→ applying seed %s", filepath.Base(seedFile))

		if err := executeSeedFile(ctx, tx, seedFile); err != nil {
			return err
		}
	}

	return nil
}

func executeSeedFile(ctx context.Context, tx *sql.Tx, path string) error {
	// #nosec G304 -- file path is intentionally provided via CLI flag
	script, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("read %q: %w", path, err)
	}

	if _, err := tx.ExecContext(ctx, string(script)); err != nil {
		return fmt.Errorf("%s: %w", filepath.Base(path), err)
	}

	return nil
}
