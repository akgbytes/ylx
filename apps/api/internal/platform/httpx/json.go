package httpx

import (
	"encoding/json"
	"errors"
	"io"
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

func DecodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(dst); err != nil {
		var typeErr *json.UnmarshalTypeError
		if errors.As(err, &typeErr) {
			WriteValidationError(w, typeErr.Field, "invalid field error")
			return false
		}

		WriteError(w, CodeMalformedJSON, "invalid request body")
		return false
	}

	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		WriteError(w, CodeMalformedJSON, "invalid request body")
		return false
	}

	return true
}
