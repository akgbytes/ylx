package router

import (
	"net/http"

	"github.com/akgbytes/ylx/internal/handler"
)

func registerAuthRoutes(mux *http.ServeMux, h *handler.AuthHandler, requireAuth func(http.Handler) http.Handler) {
	mux.Handle("GET /me", requireAuth(http.HandlerFunc(h.Me)))

	mux.HandleFunc("POST /auth/signup", h.SignUp)
	mux.HandleFunc("POST /auth/signup/verify", h.VerifySignUp)
	mux.HandleFunc("POST /auth/signup/resend", h.ResendSignUp)

	mux.HandleFunc("POST /auth/signin", h.SignIn)
	mux.HandleFunc("POST /auth/logout", h.Logout)
	mux.HandleFunc("POST /auth/refresh", h.Refresh)
}
