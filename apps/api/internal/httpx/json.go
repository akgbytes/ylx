package httpx

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type OffsetPageMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalPages int `json:"total_pages"`
	TotalCount int `json:"total_count"`
}

type CursorPageMeta struct {
	NextCursor *string `json:"next_cursor,omitempty"`
	HasMore    bool    `json:"has_more"`
	Limit      int     `json:"limit"`
}

type SuccessResponse struct {
	Data any `json:"data"`
	Meta any `json:"meta,omitempty"`
}

func DecodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(dst); err != nil {
		var typeErr *json.UnmarshalTypeError
		if errors.As(err, &typeErr) {
			WriteValidationError(w, typeErr.Field, "field has an invalid type")
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

func WriteJSON(w http.ResponseWriter, status int, data any) {
	writeJSON(w, status, SuccessResponse{
		Data: data,
	})
}

func WriteJSONWithMeta(w http.ResponseWriter, status int, data, meta any) {
	writeJSON(w, status, SuccessResponse{
		Data: data,
		Meta: meta,
	})
}

func writeJSON(w http.ResponseWriter, status int, response SuccessResponse) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}
