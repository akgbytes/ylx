package app

import (
	"database/sql"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"

	"github.com/akgbytes/ylx/internal/identity"
	"github.com/akgbytes/ylx/internal/platform/config"
	"github.com/akgbytes/ylx/internal/platform/health"
	"github.com/akgbytes/ylx/internal/platform/middleware"
)

func newHandler(cfg *config.Config, logger zerolog.Logger, db *sql.DB, rdb *redis.Client) http.Handler {
	mux := http.NewServeMux()

	health.NewHandler().RegisterRoutes(mux)

	identityModule := identity.New(identity.Deps{
		Config: cfg,
		DB:     db,
		Redis:  rdb,
	})

	identityModule.RegisterRoutes(mux)

	return middleware.RequestID(logger)(mux)
}
