package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"

	"github.com/akgbytes/ylx/internal/config"
	"github.com/akgbytes/ylx/internal/database"
	"github.com/akgbytes/ylx/internal/middleware"
	"github.com/akgbytes/ylx/internal/redis"
	"github.com/akgbytes/ylx/internal/router"
)

type Application struct {
	logger zerolog.Logger
	cfg    *config.Config
}

func NewApplication(cfg *config.Config, logger zerolog.Logger) *Application {
	return &Application{
		logger: logger,
		cfg:    cfg,
	}
}

func (app *Application) Run() error {
	dbCtx, dbCancel := context.WithTimeout(context.Background(), app.cfg.Database.DatabaseConnectTimeout)
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

	rdbCtx, rdbCancel := context.WithTimeout(context.Background(), app.cfg.Redis.RedisConnectTimeout)
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

	asynqClient := asynq.NewClientFromRedisClient(rdb)

	defer func() {
		if err := asynqClient.Close(); err != nil {
			app.logger.Error().Err(err).Msg("close Asynq client")
		}
	}()

	middlewares := middleware.NewMiddlewares(app.logger, &app.cfg.Auth, app.cfg.Server.Env == "prod")
	handler := router.NewRouter(app.cfg, db, rdb, asynqClient, middlewares)
	handler = middlewares.RequestID(handler)

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
