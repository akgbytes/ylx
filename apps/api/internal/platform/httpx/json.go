package httpx

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type Meta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalPages int `json:"total_pages"`
}

type APIResponse struct {
	Data any   `json:"data"`
	Meta *Meta `json:"meta,omitempty"`
}

func DecodeJSON(body io.Reader, dst any) error {
	decoder := json.NewDecoder(body)

	if err := decoder.Decode(dst); err != nil {
		if errors.Is(err, io.EOF) {
			return &DecodeError{Kind: DecodeEmptyBody, Err: err}
		}

		var typeErr *json.UnmarshalTypeError
		if errors.As(err, &typeErr) {
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

func WriteNoContent(w http.ResponseWriter) error {
	w.Header().Del("Content-Type")
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func WriteJSON(w http.ResponseWriter, status int, data any, meta *Meta) error {
	body, err := encode(APIResponse{
		Data: data,
		Meta: meta,
	})
	if err != nil {
		return err
	}

	return writeResponse(w, status, body)
}

func encode(value any) ([]byte, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(value); err != nil {
		return nil, err
	}

	return body.Bytes(), nil
}

func writeResponse(w http.ResponseWriter, status int, body []byte) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	n, err := w.Write(body)
	if err != nil {
		return err
	}
	if n != len(body) {
		return io.ErrShortWrite
	}

	return nil
}
