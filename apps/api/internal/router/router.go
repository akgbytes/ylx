package router

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/akgbytes/ylx/internal/handler"
)

func NewRouter(db *sql.DB, logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()
	handlers := handler.NewHandlers(db, logger)

	registerHealthRoutes(mux, handlers.Health)
	registerListingsRoutes(mux, handlers.Listings)

	return mux
}
