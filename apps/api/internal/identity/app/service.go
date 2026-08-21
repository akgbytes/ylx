package app

import (
	"context"
	"uuid"

	"github.com/akgbytes/ylx/internal/identity/adapters/otpstore"
	"github.com/akgbytes/ylx/internal/identity/domain"
	"github.com/akgbytes/ylx/internal/platform/config"
)

type Service struct {
	users      UserStore
	cfg        config.AuthConfig
	challenges *otpstore.Store
}

type UserStore interface {
	EmailTaken(ctx context.Context, email string) (bool, error)
	Create(ctx context.Context, user domain.User) (domain.User, error)
	ByEmail(ctx context.Context, email string) (domain.User, error)
	ByID(ctx context.Context, userID uuid.UUID) (domain.User, error)
}

type Deps struct {
	Users      UserStore
	Config     config.AuthConfig
	Challenges *otpstore.Store
}

func NewService(deps Deps) *Service {
	return &Service{
		users:      deps.Users,
		cfg:        deps.Config,
		challenges: deps.Challenges,
	}
}
