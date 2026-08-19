package identity

import (
	"net/http"

	"github.com/akgbytes/ylx/internal/identity/ports/rest"
)

type Module struct {
	handler *rest.Handler
}

func New() *Module {
	return &Module{
		handler: rest.NewHandler(),
	}
}

func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	m.handler.RegisterRoutes(mux)
}
