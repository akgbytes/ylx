package config

import (
	"errors"
	"os"
	"strings"
	"time"
)

type ServerConfig struct {
	Addr              string
	Env               string
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ShutdownTimeout   time.Duration
}

func loadServerConfig() (ServerConfig, error) {
	readHeaderTimeout, err := parseDuration("READ_HEADER_TIMEOUT")
	if err != nil {
		return ServerConfig{}, err
	}

	readTimeout, err := parseDuration("READ_TIMEOUT")
	if err != nil {
		return ServerConfig{}, err
	}

	writeTimeout, err := parseDuration("WRITE_TIMEOUT")
	if err != nil {
		return ServerConfig{}, err
	}

	idleTimeout, err := parseDuration("IDLE_TIMEOUT")
	if err != nil {
		return ServerConfig{}, err
	}

	shutdownTimeout, err := parseDuration("SHUTDOWN_TIMEOUT")
	if err != nil {
		return ServerConfig{}, err
	}

	return ServerConfig{
		Addr:              os.Getenv("ADDR"),
		Env:               os.Getenv("ENV"),
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
		ShutdownTimeout:   shutdownTimeout,
	}, nil
}

func (c *ServerConfig) validate() error {
	if c.Addr = strings.TrimSpace(c.Addr); c.Addr == "" {
		return errors.New("invalid configuration: ADDR is required")
	}

	if c.Env = strings.TrimSpace(c.Env); c.Env == "" {
		return errors.New("invalid configuration: ENV is required")
	}

	if c.Env != "dev" && c.Env != "prod" {
		return errors.New("invalid configuration: ENV must be 'dev' or 'prod'")
	}

	if c.ReadHeaderTimeout <= 0 {
		return errors.New("invalid configuration: READ_HEADER_TIMEOUT must be greater than 0")
	}

	if c.ReadTimeout <= 0 {
		return errors.New("invalid configuration: READ_TIMEOUT must be greater than 0")
	}

	if c.WriteTimeout <= 0 {
		return errors.New("invalid configuration: WRITE_TIMEOUT must be greater than 0")
	}

	if c.IdleTimeout <= 0 {
		return errors.New("invalid configuration: IDLE_TIMEOUT must be greater than 0")
	}

	if c.ShutdownTimeout <= 0 {
		return errors.New("invalid configuration: SHUTDOWN_TIMEOUT must be greater than 0")
	}

	return nil
}
