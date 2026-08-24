package rest

import (
	"net/http"

	"github.com/akgbytes/ylx/internal/identity/app"
	"github.com/akgbytes/ylx/internal/platform/httpx"
)

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	var req signupRequest

	if err := httpx.DecodeJSON(r.Body, &req); err != nil {
		httpx.WriteDecodeError(w, err)
		return
	}

	req.normalize()

	if field, err := req.validate(); err != nil {
		httpx.WriteValidationError(w, field, err.Error())
		return
	}

	started, err := h.service.StartSignup(
		r.Context(),
		app.SignupInput{
			Name:     req.Name,
			Email:    req.Email,
			Password: req.Password,
		},
	)
	if err != nil {
		h.writeError(w, r, "start signup", err)
		return
	}

	httpx.WriteJSON(w, http.StatusAccepted, signupResponse{RetryAt: started.RetryAt})
}

func (h *Handler) ResendSignup(w http.ResponseWriter, r *http.Request) {
	var req resendSignupRequest

	if err := httpx.DecodeJSON(r.Body, &req); err != nil {
		httpx.WriteDecodeError(w, err)
		return
	}

	req.normalize()

	if field, err := req.validate(); err != nil {
		httpx.WriteValidationError(w, field, err.Error())
		return
	}

	started, err := h.service.ResendSignup(r.Context(), req.Email)
	if err != nil {
		h.writeError(w, r, "resend signup otp", err)
		return
	}

	httpx.WriteJSON(w, http.StatusAccepted, signupResponse{RetryAt: started.RetryAt})
}

func (h *Handler) VerifySignup(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, TempResponse{Message: "verifying email..."})
}
