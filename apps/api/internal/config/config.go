package config

import (
	"log"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Addr         string        `validate:"required"`
	Env          string        `validate:"required,oneof=dev prod"`
	ReadTimeout  time.Duration `validate:"required,gt=0"`
	WriteTimeout time.Duration `validate:"required,gt=0"`
	IdleTimeout  time.Duration `validate:"required,gt=0"`
}

func MustLoad() *Config {
	validate := validator.New()

	cfg := &Config{
		Addr:         os.Getenv("ADDR"),
		Env:          os.Getenv("ENV"),
		ReadTimeout:  mustParseDuration("READ_TIMEOUT"),
		WriteTimeout: mustParseDuration("WRITE_TIMEOUT"),
		IdleTimeout:  mustParseDuration("IDLE_TIMEOUT"),
	}

	if err := validate.Struct(cfg); err != nil {
		log.Fatal("validation failed\n", err)
	}

	return cfg
}

func mustParseDuration(key string) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("validation failed\n%s is required", key)
	}

	d, err := time.ParseDuration(value)
	if err != nil {
		log.Fatalf("validation failed\ninvalid %s: %v", key, err)
	}

	return d
}
