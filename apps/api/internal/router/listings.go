package router

import (
	"net/http"

	"github.com/akgbytes/ylx/internal/handler"
)

func registerListingRoutes(mux *http.ServeMux, h *handler.ListingsHandler) {
	mux.HandleFunc("GET /listings", h.List)
	mux.HandleFunc("DELETE /listings/{id}", h.Delete)
}
