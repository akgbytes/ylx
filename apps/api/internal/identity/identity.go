package identity

import (
	"database/sql"
	"net/http"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"

	"github.com/akgbytes/ylx/internal/identity/adapters/db"
	"github.com/akgbytes/ylx/internal/identity/adapters/otpstore"
	"github.com/akgbytes/ylx/internal/identity/adapters/task"
	"github.com/akgbytes/ylx/internal/identity/app"
	"github.com/akgbytes/ylx/internal/identity/ports/rest"
	"github.com/akgbytes/ylx/internal/platform/config"
	"github.com/akgbytes/ylx/internal/platform/mailer"
)

type Module struct {
	handler *rest.Handler
}

type Deps struct {
	Config      *config.Config
	DB          *sql.DB
	Redis       *redis.Client
	AsynqClient *asynq.Client
}

func New(deps Deps) *Module {
	service := app.NewService(
		app.Deps{
			Users:      db.NewUserStore(deps.DB),
			Config:     deps.Config.Auth,
			Challenges: otpstore.NewStore(deps.Redis, deps.Config.Auth),
			Dispatcher: task.NewDispatcher(deps.AsynqClient),
		},
	)

	return &Module{
		handler: rest.NewHandler(service),
	}
}

func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	m.handler.RegisterRoutes(mux)
}

type Worker struct {
	handler *task.Handler
}

type WorkerDeps struct {
	Config *config.Config
	Redis  *redis.Client
	Sender mailer.Sender
	Logger zerolog.Logger
}

const Queue = task.Queue

func NewWorker(deps WorkerDeps) *Worker {
	return &Worker{
		handler: task.NewHandler(deps.Sender, otpstore.NewStore(deps.Redis, deps.Config.Auth), deps.Logger),
	}
}

func (w *Worker) RegisterTasks(mux *asynq.ServeMux) {
	w.handler.Register(mux)
}
