package bootstrap

import (
	"os"

	"github.com/rs/zerolog"

	"github.com/akgbytes/ylx/internal/config"
	"github.com/akgbytes/ylx/internal/logger"
)

type Runtime struct {
	Config *config.Config
	Logger zerolog.Logger
}

const timeFormat = "2006-01-02 15:04:05"

func NewBootstrapLogger() zerolog.Logger {
	return zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: timeFormat}).With().Timestamp().Logger()
}

func Load() (*Runtime, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	logger, err := logger.New(&cfg.Log, timeFormat)
	if err != nil {
		return nil, err
	}

	return &Runtime{
		Config: cfg,
		Logger: logger,
	}, nil
}
