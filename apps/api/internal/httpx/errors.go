package httpx

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type ErrorMeta struct {
	RetryAfter *int64     `json:"retry_after_seconds,omitempty"`
	RetryAt    *time.Time `json:"retry_after_at,omitempty"`
	Limit      *int       `json:"limit,omitempty"`
	Remaining  *int       `json:"remaining,omitempty"`
	ResetAt    *time.Time `json:"reset_at,omitempty"`
	Field      *string    `json:"field,omitempty"`
}

type APIError struct {
	Code    ErrorCode  `json:"code"`
	Message string     `json:"message"`
	Meta    *ErrorMeta `json:"meta,omitempty"`
}

type ErrorResponse struct {
	Error APIError `json:"error"`
}

func WriteCooldownError(w http.ResponseWriter, code ErrorCode, message string, retryAt time.Time) {
	secs := max(int64(time.Until(retryAt).Seconds()), 0)
	w.Header().Set("Retry-After", strconv.FormatInt(secs, 10))
	writeError(w, APIError{
		Code:    code,
		Message: message,
		Meta: &ErrorMeta{
			RetryAfter: &secs,
			RetryAt:    &retryAt,
		},
	})
}

func WriteRateLimitError(w http.ResponseWriter, code ErrorCode, message string, limit, remaining int, resetAt time.Time) {
	secs := max(int64(time.Until(resetAt).Seconds()), 0)
	w.Header().Set("Retry-After", strconv.FormatInt(secs, 10))
	writeError(w, APIError{
		Code:    code,
		Message: message,
		Meta: &ErrorMeta{
			Limit:     &limit,
			Remaining: &remaining,
			ResetAt:   &resetAt,
		},
	})
}

func WriteValidationError(w http.ResponseWriter, field, message string) {
	writeError(w, APIError{
		Code:    CodeValidation,
		Message: message,
		Meta:    &ErrorMeta{Field: &field},
	})
}

func WriteError(w http.ResponseWriter, code ErrorCode, message string) {
	writeError(w, APIError{Code: code, Message: message})
}

func writeError(w http.ResponseWriter, apiError APIError) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusFor(apiError.Code))

	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Error: apiError,
	})
}
