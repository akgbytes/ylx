package listing

import (
	"database/sql"
	"net/http"

	"github.com/akgbytes/ylx/internal/listing/adapters/db"
	"github.com/akgbytes/ylx/internal/listing/app"
	rest "github.com/akgbytes/ylx/internal/listing/ports"
)

type Module struct {
	handler *rest.Handler
}

type Deps struct {
	DB *sql.DB
}

func New(deps Deps) *Module {
	service := app.NewService(db.NewStore(deps.DB))

	return &Module{handler: rest.NewHandler(service)}
}

func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	m.handler.RegisterRoutes(mux)
}
