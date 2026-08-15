package config

import (
	"errors"
	"os"
	"strings"
)

type LogConfig struct {
	Level  string
	Format string
}

func loadLogConfig() LogConfig {
	return LogConfig{
		Level:  os.Getenv("LOG_LEVEL"),
		Format: os.Getenv("LOG_FORMAT"),
	}
}

func (c *LogConfig) validate() error {
	if c.Level = strings.TrimSpace(c.Level); c.Level == "" {
		return errors.New("invalid configuration: LOG_LEVEL is required")
	}

	if c.Level != "debug" && c.Level != "info" && c.Level != "warn" && c.Level != "error" && c.Level != "fatal" {
		return errors.New(
			"invalid configuration: LOG_LEVEL must be 'debug', 'info', 'warn', 'error', or 'fatal'",
		)
	}

	if c.Format = strings.TrimSpace(c.Format); c.Format == "" {
		return errors.New("invalid configuration: LOG_FORMAT is required")
	}

	if c.Format != "console" && c.Format != "json" {
		return errors.New("invalid configuration: LOG_FORMAT must be 'console' or 'json'")
	}

	return nil
}
