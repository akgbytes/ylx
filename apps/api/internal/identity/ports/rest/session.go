package rest

import (
	"net/http"

	"github.com/rs/zerolog"

	"github.com/akgbytes/ylx/internal/platform/httpx"
)

type TempResponse struct {
	Message string `json:"message"`
}

func (h *Handler) Signin(w http.ResponseWriter, r *http.Request) {
	var req signinRequest

	if err := httpx.DecodeJSON(r.Body, &req); err != nil {
		httpx.WriteDecodeError(w, err)
		return
	}

	req.normalize()

	if field, err := req.validate(); err != nil {
		httpx.WriteValidationError(w, field, err.Error())
		return
	}

	user, tokens, err := h.service.SignIn(r.Context(), req.Email, req.Password)
	if err != nil {
		h.writeError(w, r, "sign in", err)
		return
	}

	zerolog.Ctx(r.Context()).Info().Str("user_id", user.ID.String()).Msg("user signed in")

	h.cookies.Set(w, tokens)
	httpx.WriteJSON(w, http.StatusOK, newUserResponse(user))
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, TempResponse{Message: "logging out..."})
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, TempResponse{Message: "refreshing tokens..."})
}
