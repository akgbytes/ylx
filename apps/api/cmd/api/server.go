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

type apiServer struct {
	logger                 *slog.Logger
	config                 *config.Config
	databaseConnectTimeout time.Duration
}

func newAPIServer(logger *slog.Logger, config *config.Config, databaseConnectTimeout time.Duration) *apiServer {
	return &apiServer{
		logger:                 logger,
		config:                 config,
		databaseConnectTimeout: databaseConnectTimeout,
	}
}

func (server *apiServer) run() error {
	connectCtx, cancel := context.WithTimeout(context.Background(), server.databaseConnectTimeout)
	db, err := database.Connect(connectCtx, *server.config)
	cancel()

	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			server.logger.Error("close database", "error", err)
		}
	}()

	server.logger.Info("database connected")

	handler := router.NewRouter(db, server.logger)

	httpServer := http.Server{
		Addr:         server.config.Addr,
		Handler:      handler,
		ReadTimeout:  server.config.ReadTimeout,
		WriteTimeout: server.config.WriteTimeout,
		IdleTimeout:  server.config.IdleTimeout,
	}

	server.logger.Info("server listening", "addr", httpServer.Addr)

	return httpServer.ListenAndServe()
}
