package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/akgbytes/ylx/internal/config"
	"github.com/akgbytes/ylx/internal/database"
	"github.com/akgbytes/ylx/internal/router"
)

type application struct {
	logger           *slog.Logger
	cfg              *config.Config
	dbStartupTimeout time.Duration
}

func newApp(logger *slog.Logger, cfg *config.Config, dbStartupTimeout time.Duration) *application {
	return &application{
		logger:           logger,
		cfg:              cfg,
		dbStartupTimeout: dbStartupTimeout,
	}
}

func (app *application) Run() error {
	app.logger.Info("initializing application")

	dbCtx, dbCancel := context.WithTimeout(context.Background(), app.dbStartupTimeout)
	db, err := database.Connect(dbCtx, *app.cfg)
	dbCancel()

	if err != nil {
		return fmt.Errorf("init database: %w", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			app.logger.Error("failed to close database", "error", err)
		}
	}()

	app.logger.Info("database connected")

	handler := router.New(db, app.logger)

	server := http.Server{
		Addr:         app.cfg.Addr,
		Handler:      handler,
		ReadTimeout:  app.cfg.ReadTimeout,
		WriteTimeout: app.cfg.WriteTimeout,
		IdleTimeout:  app.cfg.IdleTimeout,
	}

	app.logger.Info("starting ylx api server", "addr", server.Addr)

	return server.ListenAndServe()
}
