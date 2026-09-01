package app

import (
	"context"
	"time"
	"uuid"

	"github.com/akgbytes/ylx/internal/identity/internal/adapters/otpstore"
	"github.com/akgbytes/ylx/internal/identity/internal/adapters/token"
	"github.com/akgbytes/ylx/internal/identity/internal/domain"
	"github.com/akgbytes/ylx/internal/platform/config"
)

type UserStore interface {
	EmailTaken(ctx context.Context, email string) (bool, error)
	Create(ctx context.Context, user domain.User) (domain.User, error)
	ByEmail(ctx context.Context, email string) (domain.User, error)
	ByID(ctx context.Context, userID uuid.UUID) (domain.User, error)
}

type SignupOTP struct {
	Recipient string
	Email     string
	EmailHash string
	OTP       string
	OTPHash   string
	ExpiresAt time.Time
}

type OTPDispatcher interface {
	DispatchSignupOTP(ctx context.Context, payload SignupOTP) error
}

type Service struct {
	users      UserStore
	cfg        config.AuthConfig
	challenges *otpstore.Store
	dispatcher OTPDispatcher
	signer     *token.Signer
}

type Deps struct {
	Users      UserStore
	Config     config.AuthConfig
	Challenges *otpstore.Store
	Dispatcher OTPDispatcher
	Signer     *token.Signer
}

func NewService(deps Deps) *Service {
	return &Service{
		users:      deps.Users,
		cfg:        deps.Config,
		challenges: deps.Challenges,
		dispatcher: deps.Dispatcher,
		signer:     deps.Signer,
	}
}
