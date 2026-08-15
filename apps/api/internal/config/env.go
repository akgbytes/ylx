package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

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
