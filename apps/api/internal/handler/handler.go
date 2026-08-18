package handler

import (
	"database/sql"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"

	"github.com/akgbytes/ylx/internal/config"
)

type Handlers struct {
	Health  *HealthHandler
	Listing *ListingHandler
	Auth    *AuthHandler
}

func NewHandlers(cfg *config.Config, db *sql.DB, rdb *redis.Client, asynqClient *asynq.Client) *Handlers {
	return &Handlers{
		Health:  NewHealthHandler(),
		Listing: NewListingHandler(db),
		Auth:    NewAuthHandler(cfg, db, rdb, asynqClient),
	}
}
