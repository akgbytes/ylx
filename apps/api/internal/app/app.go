package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rs/zerolog"

	"github.com/akgbytes/ylx/internal/config"
	"github.com/akgbytes/ylx/internal/database"
	"github.com/akgbytes/ylx/internal/middleware"
	"github.com/akgbytes/ylx/internal/router"
)

type Application struct {
	logger zerolog.Logger
	config *config.Config
}

func NewApplication(config *config.Config, logger zerolog.Logger) *Application {
	return &Application{
		logger: logger,
		config: config,
	}
}

func (app *Application) Run() error {
	connectCtx, cancel := context.WithTimeout(context.Background(), app.config.Server.DatabaseConnectTimeout)
	db, err := database.Connect(connectCtx, app.config.Database)
	cancel()

	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			app.logger.Err(err).Msg("close database")
		}
	}()

	app.logger.Info().Msg("database connected")

	handler := router.NewRouter(db)
	handler = middleware.RequestID(app.logger)(handler)

	httpServer := http.Server{
		Addr:              app.config.Server.Addr,
		Handler:           handler,
		ReadTimeout:       app.config.Server.ReadTimeout,
		ReadHeaderTimeout: app.config.Server.ReadHeaderTimeout,
		WriteTimeout:      app.config.Server.WriteTimeout,
		IdleTimeout:       app.config.Server.IdleTimeout,
	}

	app.logger.Info().Str("addr", httpServer.Addr).Msg("server listening")

	return httpServer.ListenAndServe()
}
