package health

import (
	"encoding/json"
	"net/http"
)

type Handler struct{}

type Response struct {
	Status string `json:"status"`
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(Response{
		Status: "ok",
	})
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /healthz", h.Healthz)
}
