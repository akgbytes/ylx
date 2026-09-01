package rest

import (
	"net/http"

	"github.com/rs/zerolog"

	"github.com/akgbytes/ylx/internal/listing/internal/app"
	"github.com/akgbytes/ylx/internal/platform/httpx"
)

type Handler struct {
	service *app.Service
}

func NewHandler(service *app.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /listings", h.List)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	listings, err := h.service.List(r.Context())
	if err != nil {
		h.writeError(w, r, "list listings", err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, newListingResponses(listings))
}

func (h *Handler) writeError(w http.ResponseWriter, r *http.Request, operation string, err error) {
	zerolog.Ctx(r.Context()).Err(err).Msg(operation)
	httpx.WriteError(w, httpx.CodeInternal, "internal server error")
}
