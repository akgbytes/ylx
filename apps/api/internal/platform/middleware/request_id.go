package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func RequestID(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := uuid.NewString()
			requestLogger := logger.With().Str("request_id", requestID).Logger()
			ctx := requestLogger.WithContext(r.Context())

			w.Header().Set("X-Request-ID", requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
