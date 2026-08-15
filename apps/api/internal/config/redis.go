package config

import (
	"errors"
	"os"
	"strings"
	"time"
)

type RedisConfig struct {
	URL                 string
	RedisConnectTimeout time.Duration
}

func loadRedisConfig() (RedisConfig, error) {
	redisConnectTimeout, err := parseDuration("REDIS_CONNECT_TIMEOUT")
	if err != nil {
		return RedisConfig{}, err
	}

	return RedisConfig{
		URL:                 os.Getenv("REDIS_URL"),
		RedisConnectTimeout: redisConnectTimeout,
	}, nil
}

func (c *RedisConfig) validate() error {
	if c.URL = strings.TrimSpace(c.URL); c.URL == "" {
		return errors.New("invalid configuration: REDIS_URL is required")
	}

	if c.RedisConnectTimeout <= 0 {
		return errors.New("invalid configuration: REDIS_CONNECT_TIMEOUT must be greater than 0")
	}

	return nil
}
