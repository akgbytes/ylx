package rest

import (
	"net/http"

	"github.com/akgbytes/ylx/internal/identity/app"
	"github.com/akgbytes/ylx/internal/platform/httpx"
)

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	var payload signupRequest

	if err := httpx.DecodeJSON(r.Body, &payload); err != nil {
		httpx.WriteDecodeError(w, err)
		return
	}

	payload.normalize()

	if field, err := payload.validate(); err != nil {
		httpx.WriteValidationError(w, field, err.Error())
		return
	}

	started, err := h.service.StartSignup(
		r.Context(),
		app.SignupInput{
			Name:     payload.Name,
			Email:    payload.Email,
			Password: payload.Password,
		},
	)
	if err != nil {
		h.writeError(w, r, "start signup", err)
		return
	}

	httpx.WriteJSON(w, http.StatusAccepted, signupResponse{RetryAt: started.RetryAt})
}

func (h *Handler) ResendSignup(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, TempResponse{Message: "sending code again..."})
}

func (h *Handler) VerifySignup(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, TempResponse{Message: "verifying email..."})
}
