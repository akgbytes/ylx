package config

import (
	"errors"
	"os"
	"strings"
	"time"
)

type RedisConfig struct {
	URL            string
	ConnectTimeout time.Duration
}

func loadRedisConfig() (RedisConfig, error) {
	connectTimeout, err := parseDuration("REDIS_CONNECT_TIMEOUT")
	if err != nil {
		return RedisConfig{}, err
	}

	return RedisConfig{
		URL:            os.Getenv("REDIS_URL"),
		ConnectTimeout: connectTimeout,
	}, nil
}

func (c *RedisConfig) validate() error {
	if c.URL = strings.TrimSpace(c.URL); c.URL == "" {
		return errors.New("invalid configuration: REDIS_URL is required")
	}

	if c.ConnectTimeout <= 0 {
		return errors.New("invalid configuration: REDIS_CONNECT_TIMEOUT must be greater than 0")
	}

	return nil
}
