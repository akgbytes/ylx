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
	defaultSeedDir = "seeds"
	defaultTimeout = 30 * time.Second
)

func main() {
	if err := run(); err != nil {
		color.Red("✘ %v", err)
		os.Exit(1)
	}
}

func run() error {
	dir := flag.String("dir", defaultSeedDir, "directory containing seed files")
	file := flag.String("file", "", "seed file to execute")
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

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	// Safe even after Commit(). Rollback() will simply return sql.ErrTxDone.
	defer func() {
		_ = tx.Rollback()
	}()

	if err := runSeed(ctx, tx, *dir, *file); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	color.Green("✔ seeding completed successfully")

	return nil
}

func runSeed(ctx context.Context, tx *sql.Tx, dir, file string) error {
	if file == "" {
		return runDirectory(ctx, tx, dir)
	}

	if filepath.Ext(file) != ".sql" {
		return errors.New("seed file must have a .sql extension")
	}

	path := file

	// Resolve relative paths against the configured seed directory.
	// Absolute paths are used as-is.
	if !filepath.IsAbs(path) {
		path = filepath.Join(dir, path)
	}

	return runFile(ctx, tx, path)
}

func runDirectory(ctx context.Context, tx *sql.Tx, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read seed directory: %w", err)
	}

	var files []string

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".sql" {
			continue
		}

		files = append(files, filepath.Join(dir, entry.Name()))
	}

	if len(files) == 0 {
		return fmt.Errorf("no .sql seed files found in %q", dir)
	}

	// Execute in alphabetical order (001_, 002_, 003_, ...)
	sort.Strings(files)

	for _, file := range files {
		color.Cyan("→ seeding %s", filepath.Base(file))

		if err := runFile(ctx, tx, file); err != nil {
			return err
		}
	}

	return nil
}

func runFile(ctx context.Context, tx *sql.Tx, path string) error {
	// #nosec G304 -- file path is intentionally provided via CLI flag
	query, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("read %q: %w", path, err)
	}

	if _, err := tx.ExecContext(ctx, string(query)); err != nil {
		return fmt.Errorf("%s: %w", filepath.Base(path), err)
	}

	return nil
}
