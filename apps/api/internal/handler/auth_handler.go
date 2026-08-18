package handler

import (
	"database/sql"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"

	"github.com/akgbytes/ylx/internal/auth"
	"github.com/akgbytes/ylx/internal/config"
)

type AuthHandler struct {
	db          *sql.DB
	rdb         *redis.Client
	cfg         *config.AuthConfig
	asynqClient *asynq.Client
	sessions    *auth.SessionManager
}

func NewAuthHandler(
	cfg *config.Config,
	db *sql.DB,
	rdb *redis.Client,
	asynqClient *asynq.Client,
) *AuthHandler {
	return &AuthHandler{
		cfg:         &cfg.Auth,
		db:          db,
		rdb:         rdb,
		asynqClient: asynqClient,
		sessions:    auth.NewSessionManager(&cfg.Auth, cfg.Server.Env == "prod"),
	}
}
