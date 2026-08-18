package middleware

import (
	"net/http"

	"github.com/rs/zerolog"

	"github.com/akgbytes/ylx/internal/auth"
	"github.com/akgbytes/ylx/internal/config"
)

type Middlewares struct {
	RequestID   func(http.Handler) http.Handler
	RequireAuth func(http.Handler) http.Handler
}

func NewMiddlewares(logger zerolog.Logger, cfg *config.AuthConfig, secure bool) *Middlewares {
	return &Middlewares{
		RequestID:   RequestID(logger),
		RequireAuth: RequireAuth(cfg, auth.NewSessionManager(cfg, secure)),
	}
}
