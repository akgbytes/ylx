package rest

import (
	"errors"
	"net/http"

	"github.com/rs/zerolog"

	"github.com/akgbytes/ylx/internal/identity/app"
	"github.com/akgbytes/ylx/internal/identity/domain"
	"github.com/akgbytes/ylx/internal/platform/httpx"
)

type Handler struct {
	service *app.Service
}

func NewHandler(service *app.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/signup", h.Signup)
	mux.HandleFunc("POST /auth/signup/verify", h.VerifySignup)
	mux.HandleFunc("POST /auth/signup/resend", h.ResendSignup)

	mux.HandleFunc("POST /auth/signin", h.Signin)
	mux.HandleFunc("POST /auth/logout", h.Logout)
	mux.HandleFunc("POST /auth/refresh", h.Refresh)
}

func (h *Handler) writeError(w http.ResponseWriter, r *http.Request, operation string, err error) {
	logger := zerolog.Ctx(r.Context())

	if cooldown, ok := errors.AsType[*domain.CooldownError](err); ok {
		httpx.WriteRateLimitError(
			w,
			"please wait before requesting another OTP",
			cooldown.RetryAt,
		)
		return
	}

	if sendLimit, ok := errors.AsType[*domain.SendLimitError](err); ok {
		httpx.WriteRateLimitError(
			w,
			"too many verification code requests",
			sendLimit.RetryAt,
		)
		return
	}

	switch {
	case errors.Is(err, domain.ErrEmailTaken):
		httpx.WriteError(w, httpx.CodeConflict, "email is already in use")

	case errors.Is(err, domain.ErrInvalidCredentials):
		httpx.WriteError(w, httpx.CodeUnauthorized, "invalid email or password")

	case errors.Is(err, domain.ErrChallengeExpired), errors.Is(err, domain.ErrChallengeMismatch):
		httpx.WriteError(w, httpx.CodeBadRequest, "invalid verification code")

	default:
		logger.Err(err).Msg(operation)
		httpx.WriteError(w, httpx.CodeInternal, "internal server error")
	}
}
