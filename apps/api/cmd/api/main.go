package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/akgbytes/ylx/internal/config"
)

const dbStartupTimeout = 30 * time.Second

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}))

	cfg := config.MustLoad()

	app := newApp(logger, cfg, dbStartupTimeout)

	if err := app.Run(); err != nil {
		logger.Error("failed to run application", "error", err)
		os.Exit(1)
	}
}
