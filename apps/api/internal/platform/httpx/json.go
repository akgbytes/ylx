package httpx

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Data any `json:"data"`
}

func WriteNoContent(w http.ResponseWriter) {
	w.Header().Del("Content-Type")
	w.WriteHeader(http.StatusNoContent)
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	body := encode(APIResponse{
		Data: data,
	})

	writeResponse(w, status, body)
}

func encode(value any) []byte {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(value); err != nil {
		// TODO: Should I log error or return it? Abhi chor deta hu :)
		return nil
	}
	return body.Bytes()
}

func writeResponse(w http.ResponseWriter, status int, body []byte) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	// TODO: Same error doubt here
	_, _ = w.Write(body)
}
