package router

import (
	"net/http"

	"github.com/akgbytes/ylx/internal/handler"
)

func registerListingRoutes(mux *http.ServeMux, listingsHandler *handler.ListingHandler) {
	mux.HandleFunc("GET /listings", listingsHandler.List)
	mux.HandleFunc("POST /listings", listingsHandler.Create)
	mux.HandleFunc("DELETE /listings/{id}", listingsHandler.Delete)
}
