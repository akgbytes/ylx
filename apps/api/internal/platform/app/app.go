package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"

	"github.com/akgbytes/ylx/internal/platform/config"
	"github.com/akgbytes/ylx/internal/platform/database"
	"github.com/akgbytes/ylx/internal/platform/redis"
)

type Application struct {
	logger zerolog.Logger
	cfg    *config.Config
}

func NewApplication(config *config.Config, logger zerolog.Logger) *Application {
	return &Application{
		logger: logger,
		cfg:    config,
	}
}

func (app *Application) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	dbCtx, dbCancel := context.WithTimeout(ctx, app.cfg.Database.ConnectTimeout)
	db, err := database.Connect(dbCtx, app.cfg.Database)
	dbCancel()

	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			app.logger.Err(err).Msg("close database")
		}
	}()

	app.logger.Info().Msg("database connected")

	rdbCtx, rdbCancel := context.WithTimeout(ctx, app.cfg.Redis.ConnectTimeout)
	rdb, err := redis.NewClient(rdbCtx, app.cfg.Redis)
	rdbCancel()

	if err != nil {
		return fmt.Errorf("connect redis: %w", err)
	}

	defer func() {
		if err := rdb.Close(); err != nil {
			app.logger.Err(err).Msg("close redis")
		}
	}()

	app.logger.Info().Msg("redis connected")

	httpServer := http.Server{
		Addr:              app.cfg.Server.Addr,
		Handler:           newHandler(app.logger),
		ReadTimeout:       app.cfg.Server.ReadTimeout,
		ReadHeaderTimeout: app.cfg.Server.ReadHeaderTimeout,
		WriteTimeout:      app.cfg.Server.WriteTimeout,
		IdleTimeout:       app.cfg.Server.IdleTimeout,
	}

	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- httpServer.ListenAndServe()
	}()

	app.logger.Info().Str("addr", httpServer.Addr).Msg("server listening")

	select {
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("serve http: %w", err)
		}
		return nil

	case <-ctx.Done():
		// Stop intercepting signals so a second 'CTRL+C' kills the process outright
		// instead of waiting on a in-flight shutdown
		stop()
		app.logger.Info().Msg("shutdown signal received")
	}

	return app.shutdown(&httpServer)
}

func (app *Application) shutdown(server *http.Server) error {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), app.cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		if closeErr := server.Close(); closeErr != nil {
			return errors.Join(
				fmt.Errorf("shutdown server: %w", err),
				fmt.Errorf("close server: %w", closeErr),
			)
		}
		return fmt.Errorf("shutdown server: %w", err)
	}

	app.logger.Info().Msg("server stopped")
	return nil
}
