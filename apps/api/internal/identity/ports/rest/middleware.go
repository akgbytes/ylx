package rest

import (
	"context"
	"net/http"

	"github.com/akgbytes/ylx/internal/platform/httpx"
)

type ctxKey struct{}

var userIDKey ctxKey

func (h *Handler) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := h.cookies.AccessToken(r)
		if accessToken == "" {
			httpx.WriteError(w, httpx.CodeUnauthorized, "authentication required")
			return
		}

		claims, err := h.signer.VerifyAccess(accessToken)
		if err != nil {
			httpx.WriteError(w, httpx.CodeUnauthorized, "authentication required")
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userIDKey, claims.UserID)))
	})
}

func UserIDFrom(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	return userID, ok
}
