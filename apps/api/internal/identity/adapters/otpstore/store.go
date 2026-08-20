package otpstore

import (
	"github.com/redis/go-redis/v9"

	"github.com/akgbytes/ylx/internal/platform/config"
)

type Store struct {
	rdb *redis.Client
	cfg config.AuthConfig
}

func NewStore(rdb *redis.Client, cfg config.AuthConfig) *Store {
	return &Store{rdb: rdb, cfg: cfg}
}

func signupChallengeKey(emailHash string) string {
	return "identity:signup:challenge:" + emailHash
}

func signupVerificationAttemptsKey(emailHash string) string {
	return "identity:signup:verification-attempts:" + emailHash
}

func signupCooldownKey(emailHash string) string {
	return "identity:signup:cooldown:" + emailHash
}

func signupSendAttemptsKey(emailHash string) string {
	return "identity:signup:send-attempts:" + emailHash
}
