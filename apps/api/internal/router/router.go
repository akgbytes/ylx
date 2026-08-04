package router

import (
	"net/http"

	"github.com/akgbytes/ylx/internal/handler"
)

func New() http.Handler {
	mux := http.NewServeMux()

	health := handler.NewHealthHandler()

	mux.HandleFunc("GET /healthz", health.Healthz)

	return mux
}
