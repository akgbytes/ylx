package router

import (
	"database/sql"
	"net/http"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"

	"github.com/akgbytes/ylx/internal/config"
	"github.com/akgbytes/ylx/internal/handler"
	"github.com/akgbytes/ylx/internal/middleware"
)

func NewRouter(
	cfg *config.Config,
	db *sql.DB,
	rdb *redis.Client,
	asynqClient *asynq.Client,
	middlewares *middleware.Middlewares,
) http.Handler {
	mux := http.NewServeMux()
	handlers := handler.NewHandlers(cfg, db, rdb, asynqClient)

	registerHealthRoutes(mux, handlers.Health)
	registerListingRoutes(mux, handlers.Listing, middlewares.RequireAuth)
	registerAuthRoutes(mux, handlers.Auth, middlewares.RequireAuth)
	return mux
}
