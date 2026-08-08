package router

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/akgbytes/ylx/internal/handler"
)

func New(db *sql.DB, logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()
	handlers := handler.New(db, logger)

	registerHealthRoutes(mux, handlers.Health)
	registerListingRoutes(mux, handlers.Listing)

	return mux
}
