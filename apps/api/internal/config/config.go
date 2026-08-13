package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Log      LogConfig
}

type ServerConfig struct {
	Addr                   string
	Env                    string
	ReadTimeout            time.Duration
	WriteTimeout           time.Duration
	IdleTimeout            time.Duration
	ReadHeaderTimeout      time.Duration
	DatabaseConnectTimeout time.Duration
}

type DatabaseConfig struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxIdleTime time.Duration
	ConnMaxLifetime time.Duration
}

type LogConfig struct {
	Level  string
	Format string
}

func Load() (*Config, error) {
	readTimeout, err := parseDuration("READ_TIMEOUT")
	if err != nil {
		return nil, err
	}

	readHeaderTimeout, err := parseDuration("READ_HEADER_TIMEOUT")
	if err != nil {
		return nil, err
	}

	writeTimeout, err := parseDuration("WRITE_TIMEOUT")
	if err != nil {
		return nil, err
	}

	idleTimeout, err := parseDuration("IDLE_TIMEOUT")
	if err != nil {
		return nil, err
	}

	databaseConnectTimeout, err := parseDuration("DATABASE_CONNECT_TIMEOUT")
	if err != nil {
		return nil, err
	}

	maxOpenConns, err := parseInt("DATABASE_MAX_OPEN_CONNS")
	if err != nil {
		return nil, err
	}

	maxIdleConns, err := parseInt("DATABASE_MAX_IDLE_CONNS")
	if err != nil {
		return nil, err
	}

	connMaxIdleTime, err := parseDuration("DATABASE_CONN_MAX_IDLE_TIME")
	if err != nil {
		return nil, err
	}

	connMaxLifetime, err := parseDuration("DATABASE_CONN_MAX_LIFETIME")
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Server: ServerConfig{
			Addr:                   os.Getenv("ADDR"),
			Env:                    os.Getenv("ENV"),
			ReadTimeout:            readTimeout,
			ReadHeaderTimeout:      readHeaderTimeout,
			WriteTimeout:           writeTimeout,
			IdleTimeout:            idleTimeout,
			DatabaseConnectTimeout: databaseConnectTimeout,
		},
		Database: DatabaseConfig{
			URL:             os.Getenv("DATABASE_URL"),
			MaxOpenConns:    maxOpenConns,
			MaxIdleConns:    maxIdleConns,
			ConnMaxIdleTime: connMaxIdleTime,
			ConnMaxLifetime: connMaxLifetime,
		},
		Log: LogConfig{
			Level:  os.Getenv("LOG_LEVEL"),
			Format: os.Getenv("LOG_FORMAT"),
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if err := c.Server.validate(); err != nil {
		return err
	}

	if err := c.Database.validate(); err != nil {
		return err
	}

	if err := c.Log.validate(); err != nil {
		return err
	}

	return nil
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

	if c.DatabaseConnectTimeout <= 0 {
		return errors.New("invalid configuration: DATABASE_CONNECT_TIMEOUT must be greater than 0")
	}

	return nil
}

func (c *DatabaseConfig) validate() error {
	if c.URL = strings.TrimSpace(c.URL); c.URL == "" {
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

	return nil
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

func parseDuration(key string) (time.Duration, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return 0, fmt.Errorf("invalid configuration: %s is required", key)
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("invalid configuration: %s: %w", key, err)
	}

	return duration, nil
}

func parseInt(key string) (int, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return 0, fmt.Errorf("invalid configuration: %s is required", key)
	}

	integer, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid configuration: %s: %w", key, err)
	}
	return integer, nil
}
