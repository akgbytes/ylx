package httpx

type DecodeErrorKind string

const (
	DecodeMalformedJSON DecodeErrorKind = "malformed_json"
	DecodeTypeMismatch  DecodeErrorKind = "type_mismatch"
	DecodeEmptyBody     DecodeErrorKind = "empty_body"
	DecodeTrailingData  DecodeErrorKind = "trailing_data"
)

type DecodeError struct {
	Kind  DecodeErrorKind
	Field string
	Err   error
}

func (e *DecodeError) Error() string {
	if e == nil || e.Err == nil {
		return "decode error"
	}

	return e.Err.Error()
}
