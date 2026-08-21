package identity

import (
	"database/sql"
	"net/http"

	"github.com/redis/go-redis/v9"

	"github.com/akgbytes/ylx/internal/identity/adapters/db"
	"github.com/akgbytes/ylx/internal/identity/adapters/otpstore"
	"github.com/akgbytes/ylx/internal/identity/app"
	"github.com/akgbytes/ylx/internal/identity/ports/rest"
	"github.com/akgbytes/ylx/internal/platform/config"
)

type Module struct {
	handler *rest.Handler
}

type Deps struct {
	Config *config.Config
	DB     *sql.DB
	Redis  *redis.Client
}

func New(deps Deps) *Module {
	service := app.NewService(
		app.Deps{
			Users:      db.NewUserStore(deps.DB),
			Config:     deps.Config.Auth,
			Challenges: otpstore.NewStore(deps.Redis, deps.Config.Auth),
		},
	)

	return &Module{
		handler: rest.NewHandler(service),
	}
}

func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	m.handler.RegisterRoutes(mux)
}
