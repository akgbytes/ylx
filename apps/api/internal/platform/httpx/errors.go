package httpx

import (
	"errors"
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

func WriteError(w http.ResponseWriter, code ErrorCode, message string) error {
	body, err := encode(ErrorResponse{
		Error: APIError{
			Code:    code,
			Message: message,
		},
	})
	if err != nil {
		return err
	}

	return writeResponse(w, statusFor(code), body)
}

func WriteValidationError(w http.ResponseWriter, field, message string) error {
	body, err := encode(ErrorResponse{
		Error: APIError{
			Code:    CodeValidation,
			Message: message,
			Field:   field,
		},
	})
	if err != nil {
		return err
	}

	return writeResponse(w, statusFor(CodeValidation), body)
}

func WriteDecodeError(w http.ResponseWriter, err error) error {
	var decodeErr *DecodeError
	if !errors.As(err, &decodeErr) || decodeErr.Kind != DecodeTypeMismatch {
		return WriteError(w, CodeMalformedJSON, "invalid request body")
	}

	return WriteValidationError(w, decodeErr.Field, "invalid field type")
}
