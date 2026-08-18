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
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ReadHeaderTimeout time.Duration
}

func loadServerConfig() (ServerConfig, error) {
	readTimeout, err := parseDuration("READ_TIMEOUT")
	if err != nil {
		return ServerConfig{}, err
	}

	readHeaderTimeout, err := parseDuration("READ_HEADER_TIMEOUT")
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

	return ServerConfig{
		Addr:              os.Getenv("ADDR"),
		Env:               os.Getenv("ENV"),
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
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

	if c.ReadTimeout <= 0 {
		return errors.New("invalid configuration: READ_TIMEOUT must be greater than 0")
	}

	if c.WriteTimeout <= 0 {
		return errors.New("invalid configuration: WRITE_TIMEOUT must be greater than 0")
	}

	if c.IdleTimeout <= 0 {
		return errors.New("invalid configuration: IDLE_TIMEOUT must be greater than 0")
	}

	if c.ReadHeaderTimeout <= 0 {
		return errors.New("invalid configuration: READ_HEADER_TIMEOUT must be greater than 0")
	}

	return nil
}
