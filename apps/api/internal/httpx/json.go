package httpx

import (
	"encoding/json"
	"net/http"
)

type Meta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalPages int `json:"total_pages"`
}

type APIResponse struct {
	Data any   `json:"data"`
	Meta *Meta `json:"meta,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, data any, meta *Meta) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(APIResponse{
		Data: data,
		Meta: meta,
	})
}
