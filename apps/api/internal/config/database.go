package config

import (
	"errors"
	"os"
	"strings"
	"time"
)

type DatabaseConfig struct {
	URL                    string
	MaxOpenConns           int
	MaxIdleConns           int
	ConnMaxIdleTime        time.Duration
	ConnMaxLifetime        time.Duration
	DatabaseConnectTimeout time.Duration
}

func loadDatabaseConfig() (DatabaseConfig, error) {
	maxOpenConns, err := parseInt("DATABASE_MAX_OPEN_CONNS")
	if err != nil {
		return DatabaseConfig{}, err
	}

	maxIdleConns, err := parseInt("DATABASE_MAX_IDLE_CONNS")
	if err != nil {
		return DatabaseConfig{}, err
	}

	connMaxIdleTime, err := parseDuration("DATABASE_CONN_MAX_IDLE_TIME")
	if err != nil {
		return DatabaseConfig{}, err
	}

	connMaxLifetime, err := parseDuration("DATABASE_CONN_MAX_LIFETIME")
	if err != nil {
		return DatabaseConfig{}, err
	}

	databaseConnectTimeout, err := parseDuration("DATABASE_CONNECT_TIMEOUT")
	if err != nil {
		return DatabaseConfig{}, err
	}

	return DatabaseConfig{
		URL:                    os.Getenv("DATABASE_URL"),
		MaxOpenConns:           maxOpenConns,
		MaxIdleConns:           maxIdleConns,
		ConnMaxIdleTime:        connMaxIdleTime,
		ConnMaxLifetime:        connMaxLifetime,
		DatabaseConnectTimeout: databaseConnectTimeout,
	}, nil
}

func (c *DatabaseConfig) normalize() {
	c.URL = strings.TrimSpace(c.URL)
}

func (c *DatabaseConfig) validate() error {
	if c.URL == "" {
		return errors.New("invalid configuration: DATABASE_URL is required")
	}

	if c.MaxOpenConns <= 0 {
		return errors.New("invalid configuration: DATABASE_MAX_OPEN_CONNS must be greater than 0")
	}

	if c.MaxIdleConns <= 0 {
		return errors.New("invalid configuration: DATABASE_MAX_IDLE_CONNS must be greater than 0")
	}

	if c.MaxIdleConns > c.MaxOpenConns {
		return errors.New(
			"invalid configuration: DATABASE_MAX_IDLE_CONNS must not exceed DATABASE_MAX_OPEN_CONNS",
		)
	}

	if c.ConnMaxIdleTime <= 0 {
		return errors.New("invalid configuration: DATABASE_CONN_MAX_IDLE_TIME must be greater than 0")
	}

	if c.ConnMaxLifetime <= 0 {
		return errors.New("invalid configuration: DATABASE_CONN_MAX_LIFETIME must be greater than 0")
	}

	if c.DatabaseConnectTimeout <= 0 {
		return errors.New("invalid configuration: DATABASE_CONNECT_TIMEOUT must be greater than 0")
	}

	return nil
}
