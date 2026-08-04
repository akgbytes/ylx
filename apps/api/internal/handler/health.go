package handler

import (
	"encoding/json"
	"net/http"
)

type HealthHandler struct{}

type HealthResponse struct {
	Status string `json:"status"`
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(HealthResponse{
		Status: "ok",
	})
}
