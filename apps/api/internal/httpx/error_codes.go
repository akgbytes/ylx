package httpx

import "net/http"

type ErrorCode string

const (
	CodeMalformedJSON    ErrorCode = "malformed_json"
	CodeValidation       ErrorCode = "validation_error"
	CodeBadRequest       ErrorCode = "bad_request"
	CodeUnprocessable    ErrorCode = "unprocessable_content"
	CodeUnsupportedMedia ErrorCode = "unsupported_media_type"
	CodePayloadTooLarge  ErrorCode = "payload_too_large"

	CodeUnauthorized ErrorCode = "unauthorized"
	CodeForbidden    ErrorCode = "forbidden"

	CodeNotFound ErrorCode = "not_found"
	CodeConflict ErrorCode = "conflict"

	CodeTooManyRequests   ErrorCode = "too_many_requests"
	CodeOTPCooldownActive ErrorCode = "otp_cooldown_active"
	CodeOTPHourlyLimit    ErrorCode = "otp_hourly_limit_reached"

	CodeOTPInvalid ErrorCode = "otp_invalid"
	CodeOTPExpired ErrorCode = "otp_expired"

	CodeInternal       ErrorCode = "internal_error"
	CodeServiceUnavail ErrorCode = "service_unavailable"
)

func statusFor(code ErrorCode) int {
	switch code {
	case CodeMalformedJSON, CodeValidation, CodeBadRequest, CodeOTPInvalid:
		return http.StatusBadRequest
	case CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeForbidden:
		return http.StatusForbidden
	case CodeNotFound:
		return http.StatusNotFound
	case CodeConflict:
		return http.StatusConflict
	case CodeUnprocessable, CodeOTPExpired:
		return http.StatusUnprocessableEntity
	case CodeUnsupportedMedia:
		return http.StatusUnsupportedMediaType
	case CodePayloadTooLarge:
		return http.StatusRequestEntityTooLarge
	case CodeTooManyRequests, CodeOTPCooldownActive, CodeOTPHourlyLimit:
		return http.StatusTooManyRequests
	case CodeServiceUnavail:
		return http.StatusServiceUnavailable
	case CodeInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
