package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Addr                    string        `validate:"required"`
	Env                     string        `validate:"required,oneof=dev prod"`
	ReadTimeout             time.Duration `validate:"required,gt=0"`
	WriteTimeout            time.Duration `validate:"required,gt=0"`
	IdleTimeout             time.Duration `validate:"required,gt=0"`
	DatabaseURL             string        `validate:"required"`
	DatabaseMaxOpenConns    int           `validate:"required,gt=0"`
	DatabaseMaxIdleConns    int           `validate:"required,gt=0"`
	DatabaseConnMaxIdleTime time.Duration `validate:"required,gt=0"`
	DatabaseConnMaxLifetime time.Duration `validate:"required,gt=0"`
}

func MustLoad() *Config {
	validator := validator.New()

	cfg := &Config{
		Addr:                    os.Getenv("ADDR"),
		Env:                     os.Getenv("ENV"),
		ReadTimeout:             mustParseDuration("READ_TIMEOUT"),
		WriteTimeout:            mustParseDuration("WRITE_TIMEOUT"),
		IdleTimeout:             mustParseDuration("IDLE_TIMEOUT"),
		DatabaseURL:             os.Getenv("DATABASE_URL"),
		DatabaseMaxOpenConns:    mustParseInt("DATABASE_MAX_OPEN_CONNS"),
		DatabaseMaxIdleConns:    mustParseInt("DATABASE_MAX_IDLE_CONNS"),
		DatabaseConnMaxIdleTime: mustParseDuration("DATABASE_CONN_MAX_IDLE_TIME"),
		DatabaseConnMaxLifetime: mustParseDuration("DATABASE_CONN_MAX_LIFETIME"),
	}

	if err := validator.Struct(cfg); err != nil {
		log.Fatal("invalid configuration: ", err)
	}

	return cfg
}

func mustParseDuration(key string) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("invalid configuration: %s is required", key)
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		log.Fatalf("invalid configuration: %s: %v", key, err)
	}

	return duration
}

func mustParseInt(key string) int {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("invalid configuration: %s is required", key)
	}

	integer, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("invalid configuration: %s: %v", key, err)
	}
	return integer
}
