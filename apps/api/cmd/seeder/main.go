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

	"github.com/akgbytes/ylx/internal/config"
	"github.com/akgbytes/ylx/internal/database"
	"github.com/akgbytes/ylx/internal/logger"
)

const seedTransactionTimeout = 30 * time.Second

func main() {
	logger := logger.BootstrapLogger()

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal().Err(err).Msg("bootstrap seeder")
	}

	seedDir := flag.String("dir", "seeds", "directory containing seed files")
	seedFile := flag.String("file", "", "seed file to execute")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), seedTransactionTimeout)
	defer cancel()

	db, err := database.Connect(ctx, cfg.Database)
	if err != nil {
		logger.Fatal().Err(err).Msg("connect database")
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error().Err(err).Msg("close database")
		}
	}()

	if err := seed(ctx, db, *seedDir, *seedFile, logger); err != nil {
		logger.Fatal().Err(err).Msg("seed database")
	}

	logger.Info().Msg("database seeded")
}

func seed(ctx context.Context, db *sql.DB, seedDir, seedFile string, logger zerolog.Logger) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	files, err := filesToSeed(seedDir, seedFile)
	if err != nil {
		return err
	}

	for _, path := range files {
		logger.Info().Str("file", filepath.Base(path)).Msg("applying seed")
		if err := executeFile(ctx, tx, path); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

func filesToSeed(seedDir, seedFile string) ([]string, error) {
	if seedFile != "" {
		if filepath.Ext(seedFile) != ".sql" {
			return nil, errors.New("seed file must have a .sql extension")
		}
		if !filepath.IsAbs(seedFile) {
			seedFile = filepath.Join(seedDir, seedFile)
		}
		return []string{seedFile}, nil
	}

	entries, err := os.ReadDir(seedDir)
	if err != nil {
		return nil, fmt.Errorf("read seed directory: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".sql" {
			files = append(files, filepath.Join(seedDir, entry.Name()))
		}
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no .sql seed files found in %q", seedDir)
	}
	sort.Strings(files)
	return files, nil
}

func executeFile(ctx context.Context, tx *sql.Tx, path string) error {
	script, err := os.ReadFile(filepath.Clean(path)) // #nosec G304 -- CLI-selected seed source.
	if err != nil {
		return fmt.Errorf("read %q: %w", path, err)
	}
	if _, err := tx.ExecContext(ctx, string(script)); err != nil {
		return fmt.Errorf("execute %q: %w", filepath.Base(path), err)
	}
	return nil
}
