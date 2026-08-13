package router

import (
	"database/sql"
	"net/http"

	"github.com/akgbytes/ylx/internal/handler"
)

func NewRouter(db *sql.DB) http.Handler {
	mux := http.NewServeMux()
	handlers := handler.NewHandlers(db)

	registerHealthRoutes(mux, handlers.Health)
	registerListingsRoutes(mux, handlers.Listings)

	return mux
}
