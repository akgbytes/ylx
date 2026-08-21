package rest

import (
	"net/http"

	"github.com/akgbytes/ylx/internal/platform/httpx"
)

type TempResponse struct {
	Message string `json:"message"`
}

func (h *Handler) Signin(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, TempResponse{Message: "signing in..."})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, TempResponse{Message: "logging out..."})
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, TempResponse{Message: "refreshing tokens..."})
}
