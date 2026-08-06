package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/akgbytes/ylx/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect(ctx context.Context, cfg config.Config) (*sql.DB, error) {

	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("sql open: %w", err)
	}

	db.SetMaxOpenConns(cfg.DatabaseMaxOpenConns)
	db.SetMaxIdleConns(cfg.DatabaseMaxIdleConns)
	db.SetConnMaxIdleTime(cfg.DatabaseConnMaxIdleTime)
	db.SetConnMaxLifetime(cfg.DatabaseConnMaxLifetime)

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("database ping: %w", err)
	}

	return db, nil
}
