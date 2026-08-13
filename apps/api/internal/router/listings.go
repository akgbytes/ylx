package router

import (
	"net/http"

	"github.com/akgbytes/ylx/internal/handler"
)

func registerListingsRoutes(mux *http.ServeMux, listingsHandler *handler.ListingsHandler) {
	mux.HandleFunc("GET /listings", listingsHandler.List)
	mux.HandleFunc("DELETE /listings/{id}", listingsHandler.Delete)
}
