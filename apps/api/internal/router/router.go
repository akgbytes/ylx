package router

import (
	"database/sql"
	"net/http"

	"github.com/rs/zerolog"

	"github.com/akgbytes/ylx/internal/handler"
)

func NewRouter(db *sql.DB, logger zerolog.Logger) http.Handler {
	mux := http.NewServeMux()
	handlers := handler.NewHandlers(db, logger)

	registerHealthRoutes(mux, handlers.Health)
	registerListingsRoutes(mux, handlers.Listings)

	return mux
}
