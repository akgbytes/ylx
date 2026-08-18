package router

import (
	"net/http"

	"github.com/akgbytes/ylx/internal/handler"
)

func registerListingRoutes(mux *http.ServeMux, h *handler.ListingHandler, requireAuth func(http.Handler) http.Handler) {
	mux.HandleFunc("GET /listings", h.List)
	mux.Handle("POST /listings", requireAuth(http.HandlerFunc(h.Create)))
	mux.Handle("DELETE /listings/{id}", requireAuth(http.HandlerFunc(h.Delete)))
}
