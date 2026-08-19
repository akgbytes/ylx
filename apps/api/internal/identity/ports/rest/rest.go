package rest

import "net/http"

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/signup", h.Signup)
	mux.HandleFunc("POST /auth/signup/verify", h.VerifySignup)
	mux.HandleFunc("POST /auth/signup/resend", h.ResendSignup)

	mux.HandleFunc("POST /auth/signin", h.Signin)
	mux.HandleFunc("POST /auth/logout", h.Logout)
	mux.HandleFunc("POST /auth/refresh", h.Refresh)
}
