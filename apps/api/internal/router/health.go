package router

import (
	"net/http"

	"github.com/akgbytes/ylx/internal/handler"
)

func registerHealthRoutes(mux *http.ServeMux, h *handler.HealthHandler) {
	mux.HandleFunc("GET /healthz", h.Healthz)
}
