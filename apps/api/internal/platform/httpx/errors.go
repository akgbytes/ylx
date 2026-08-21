package httpx

import (
	"errors"
	"net/http"
	"strconv"
	"time"
)

type ErrorCode string

type ErrorMeta struct {
	Field      *string    `json:"field,omitempty"`
	RetryAfter *int64     `json:"retry_after_seconds,omitempty"`
	RetryAt    *time.Time `json:"retry_after_at,omitempty"`
}

type APIError struct {
	Code    ErrorCode  `json:"code"`
	Message string     `json:"message"`
	Meta    *ErrorMeta `json:"meta,omitempty"`
}

type ErrorResponse struct {
	Error APIError `json:"error"`
}

func WriteError(w http.ResponseWriter, code ErrorCode, message string) {
	body := encode(ErrorResponse{
		Error: APIError{
			Code:    code,
			Message: message,
		},
	})
	writeResponse(w, statusFor(code), body)
}

func WriteValidationError(w http.ResponseWriter, field, message string) {
	body := encode(ErrorResponse{
		Error: APIError{
			Code:    CodeValidation,
			Message: message,
			Meta: &ErrorMeta{
				Field: &field,
			},
		},
	})
	writeResponse(w, statusFor(CodeValidation), body)
}

func WriteRateLimitError(w http.ResponseWriter, message string, retryAt time.Time) {
	secs := max(int64(time.Until(retryAt).Seconds()), 0)
	w.Header().Set("Retry-After", strconv.FormatInt(secs, 10))

	body := encode(ErrorResponse{
		Error: APIError{
			Code:    CodeTooManyRequests,
			Message: message,
			Meta: &ErrorMeta{
				RetryAfter: &secs,
				RetryAt:    &retryAt,
			},
		},
	})
	writeResponse(w, statusFor(CodeTooManyRequests), body)
}

func WriteDecodeError(w http.ResponseWriter, err error) {
	var decodeErr *DecodeError
	if !errors.As(err, &decodeErr) || decodeErr.Kind != DecodeTypeMismatch {
		WriteError(w, CodeMalformedJSON, "invalid request body")
	}
	WriteValidationError(w, decodeErr.Field, "invalid field type")
}
