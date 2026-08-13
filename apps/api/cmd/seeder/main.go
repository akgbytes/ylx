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

	"github.com/rs/zerolog"

	"github.com/akgbytes/ylx/internal/bootstrap"
	"github.com/akgbytes/ylx/internal/database"
)

const defaultTimeout = 30 * time.Second

func main() {
	bootstraplogger := bootstrap.NewBootstrapLogger()
	runtime, err := bootstrap.Load()
	if err != nil {
		bootstraplogger.Fatal().Err(err).Msg("bootstrap seeder")
	}

	dir := flag.String("dir", "seeds", "directory containing seed files")
	file := flag.String("file", "", "seed file to execute")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	db, err := database.Connect(ctx, runtime.Config.Database)
	if err != nil {
		runtime.Logger.Fatal().Err(err).Msg("connect database")
	}
	defer func() {
		if err := db.Close(); err != nil {
			runtime.Logger.Error().Err(err).Msg("close database")
		}
	}()

	if err := seed(ctx, db, *dir, *file, runtime.Logger); err != nil {
		runtime.Logger.Fatal().Err(err).Msg("seed database")
	}

	runtime.Logger.Info().Msg("seeding completed")
}

func seed(ctx context.Context, db *sql.DB, dir, file string, logger zerolog.Logger) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	files, err := filesToSeed(dir, file)
	if err != nil {
		return err
	}

	for _, path := range files {
		logger.Info().Str("file", filepath.Base(path)).Msg("seeding")
		if err := executeFile(ctx, tx, path); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

func filesToSeed(dir, file string) ([]string, error) {
	if file != "" {
		if filepath.Ext(file) != ".sql" {
			return nil, errors.New("seed file must have a .sql extension")
		}
		if !filepath.IsAbs(file) {
			file = filepath.Join(dir, file)
		}
		return []string{file}, nil
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read seed directory: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".sql" {
			files = append(files, filepath.Join(dir, entry.Name()))
		}
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no .sql seed files found in %q", dir)
	}
	sort.Strings(files)
	return files, nil
}

func executeFile(ctx context.Context, tx *sql.Tx, path string) error {
	query, err := os.ReadFile(filepath.Clean(path)) // #nosec G304 -- CLI-selected seed source.
	if err != nil {
		return fmt.Errorf("read %q: %w", path, err)
	}
	if _, err := tx.ExecContext(ctx, string(query)); err != nil {
		return fmt.Errorf("execute %q: %w", filepath.Base(path), err)
	}
	return nil
}
