package main

import (
	"github.com/akgbytes/ylx/internal/app"
	"github.com/akgbytes/ylx/internal/config"
	"github.com/akgbytes/ylx/internal/logger"
)

func main() {
	bootstrapLogger := logger.BootstrapLogger()

	cfg, err := config.Load()
	if err != nil {
		bootstrapLogger.Fatal().Err(err).Msg("bootstrap application")
	}

	log, err := logger.New(&cfg.Log)
	if err != nil {
		bootstrapLogger.Fatal().Err(err).Msg("bootstrap application")
	}

	application := app.NewApplication(cfg, log)

	if err := application.Run(); err != nil {
		log.Fatal().Err(err).Msg("server exited")
	}
}
