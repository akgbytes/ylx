package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akgbytes/ylx/internal/platform/config"
	"github.com/akgbytes/ylx/internal/platform/database"
	"github.com/akgbytes/ylx/internal/platform/health"
	"github.com/akgbytes/ylx/internal/platform/middleware"
	"github.com/akgbytes/ylx/internal/platform/redis"

	"github.com/rs/zerolog"
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
	dbCtx, dbCancel := context.WithTimeout(context.Background(), app.cfg.Database.ConnectTimeout)
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

	rdbCtx, rdbCancel := context.WithTimeout(context.Background(), app.cfg.Redis.ConnectTimeout)
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

	mux := http.NewServeMux()

	// Register routes
	health.NewHandler().RegisterRoutes(mux)

	// Middlewares
	handler := middleware.RequestID(app.logger)(mux)

	httpServer := http.Server{
		Addr:              app.cfg.Server.Addr,
		Handler:           handler,
		ReadTimeout:       app.cfg.Server.ReadTimeout,
		ReadHeaderTimeout: app.cfg.Server.ReadHeaderTimeout,
		WriteTimeout:      app.cfg.Server.WriteTimeout,
		IdleTimeout:       app.cfg.Server.IdleTimeout,
	}

	app.logger.Info().Str("addr", httpServer.Addr).Msg("server listening")

	return httpServer.ListenAndServe()
}
