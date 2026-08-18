package middleware

import (
	"context"
	"net/http"

	"github.com/akgbytes/ylx/internal/auth"
	"github.com/akgbytes/ylx/internal/config"
	"github.com/akgbytes/ylx/internal/httpx"
)

type ctxKey struct{}

var authUserIDKey ctxKey

func RequireAuth(cfg *config.AuthConfig, sessions *auth.SessionManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(cfg.AccessTokenName)
			if err != nil {
				httpx.WriteError(w, httpx.CodeUnauthorized, "authentication required")
				return
			}

			userID, err := sessions.VerifyAccessToken(cookie.Value)
			if err != nil {
				httpx.WriteError(w, httpx.CodeUnauthorized, "authentication required")
				return
			}

			ctx := context.WithValue(r.Context(), authUserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(authUserIDKey).(string)
	return userID, ok
}
