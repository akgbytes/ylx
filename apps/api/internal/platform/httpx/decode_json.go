package httpx

import (
	"encoding/json"
	"errors"
	"io"
)

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

func DecodeJSON(body io.Reader, dst any) error {
	decoder := json.NewDecoder(body)

	if err := decoder.Decode(dst); err != nil {
		if errors.Is(err, io.EOF) {
			return &DecodeError{Kind: DecodeEmptyBody, Err: err}
		}
		if typeErr, ok := errors.AsType[*json.UnmarshalTypeError](err); ok {
			return &DecodeError{
				Kind:  DecodeTypeMismatch,
				Field: typeErr.Field,
				Err:   err,
			}
		}
		return &DecodeError{Kind: DecodeMalformedJSON, Err: err}
	}

	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return &DecodeError{Kind: DecodeTrailingData, Err: err}
	}

	return nil
}
