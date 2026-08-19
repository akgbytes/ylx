package httpx

import "net/http"

const (
	CodeMalformedJSON    ErrorCode = "malformed_json"
	CodeValidation       ErrorCode = "validation_error"
	CodeUnauthorized     ErrorCode = "unauthorized"
	CodeForbidden        ErrorCode = "forbidden"
	CodeBadRequest       ErrorCode = "bad_request"
	CodeNotFound         ErrorCode = "not_found"
	CodeConflict         ErrorCode = "conflict"
	CodeUnprocessable    ErrorCode = "unprocessable_content"
	CodeUnsupportedMedia ErrorCode = "unsupported_media_type"
	CodeInternal         ErrorCode = "internal_error"
)

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
