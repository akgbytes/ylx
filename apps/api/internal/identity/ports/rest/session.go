package rest

import "net/http"

func (h *Handler) Signin(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("signing in..."))
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("logging out..."))
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("refreshing tokens..."))
}
