package httpx

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
