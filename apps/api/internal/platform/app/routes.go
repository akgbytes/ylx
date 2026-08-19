package app

import (
	"net/http"

	"github.com/rs/zerolog"

	"github.com/akgbytes/ylx/internal/identity"
	"github.com/akgbytes/ylx/internal/platform/health"
	"github.com/akgbytes/ylx/internal/platform/middleware"
)

func newHandler(logger zerolog.Logger) http.Handler {
	mux := http.NewServeMux()

	health.NewHandler().RegisterRoutes(mux)

	identityModule := identity.New()
	identityModule.RegisterRoutes(mux)

	return middleware.RequestID(logger)(mux)
}
