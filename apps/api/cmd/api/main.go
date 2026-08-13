package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/akgbytes/ylx/internal/config"
)

const databaseConnectTimeout = 30 * time.Second

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}))

	cfg := config.MustLoad()

	api := newAPIServer(logger, cfg, databaseConnectTimeout)

	if err := api.run(); err != nil {
		logger.Error("server exited", "error", err)
		os.Exit(1)
	}
}
