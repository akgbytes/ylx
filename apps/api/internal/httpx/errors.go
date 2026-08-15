package httpx

import (
	"encoding/json"
	"net/http"
)

type ErrorCode string

type APIError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Field   string    `json:"field,omitempty"`
}

type ErrorResponse struct {
	Error APIError `json:"error"`
}

func WriteError(w http.ResponseWriter, code ErrorCode, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusFor(code))

	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Error: APIError{
			Code:    code,
			Message: message,
		},
	})
}

func WriteValidationError(w http.ResponseWriter, field, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusFor(CodeValidation))

	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Error: APIError{
			Code:    CodeValidation,
			Message: message,
			Field:   field,
		},
	})
}

func statusFor(code ErrorCode) int {
	switch code {
	case CodeBadRequest, CodeMalformedJSON, CodeValidation:
		return http.StatusBadRequest
	case CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeForbidden:
		return http.StatusForbidden
	case CodeNotFound:
		return http.StatusNotFound
	case CodeConflict:
		return http.StatusConflict
	case CodeUnprocessable:
		return http.StatusUnprocessableEntity
	case CodeUnsupportedMedia:
		return http.StatusUnsupportedMediaType
	case CodeInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
