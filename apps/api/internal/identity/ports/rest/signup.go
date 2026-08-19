package rest

import (
	"net/http"
)

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("signing up..."))
}

func (h *Handler) ResendSignup(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("sending code again..."))
}

func (h *Handler) VerifySignup(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("verifying email..."))
}
