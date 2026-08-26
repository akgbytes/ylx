package app

import (
	"context"
	"crypto/hmac"
	"errors"
	"fmt"
	"time"
	"uuid"

	"github.com/redis/go-redis/v9"

	"github.com/akgbytes/ylx/internal/identity/adapters/crypto"
	"github.com/akgbytes/ylx/internal/identity/adapters/otpstore"
	"github.com/akgbytes/ylx/internal/identity/domain"
)

type SignupInput struct {
	Name     string
	Email    string
	Password string
}

type SignupOutput struct {
	RetryAt time.Time
}

func (s *Service) StartSignup(ctx context.Context, in SignupInput) (SignupOutput, error) {
	taken, err := s.users.EmailTaken(ctx, in.Email)
	if err != nil {
		return SignupOutput{}, err
	}

	if taken {
		return SignupOutput{}, domain.ErrEmailTaken
	}

	passwordHash, err := crypto.HashPassword(in.Password, crypto.DefaultArgonParams())
	if err != nil {
		return SignupOutput{}, fmt.Errorf("hash password: %w", err)
	}

	otp, err := crypto.GenerateOTP()
	if err != nil {
		return SignupOutput{}, fmt.Errorf("generate otp: %w", err)
	}

	otpHash := crypto.HashOTP(otp, s.cfg.OTPSecretKey)
	emailHash := crypto.HashEmail(in.Email)

	reservation, err := s.challenges.Reserve(ctx, otpstore.Challenge{
		Name:         in.Name,
		Email:        in.Email,
		EmailHash:    emailHash,
		PasswordHash: passwordHash,
		OTPHash:      otpHash,
	})
	if err != nil {
		return SignupOutput{}, err
	}
	if !reservation.Allowed {
		return SignupOutput{}, reservationError(reservation)
	}

	if err := s.dispatcher.DispatchSignupOTP(ctx, SignupOTP{
		Recipient: in.Name,
		Email:     in.Email,
		EmailHash: emailHash,
		OTP:       otp,
		OTPHash:   otpHash,
		ExpiresAt: reservation.ExpiresAt,
	}); err != nil {
		if releaseErr := s.challenges.Release(ctx, emailHash, otpHash); releaseErr != nil {
			return SignupOutput{}, errors.Join(err, releaseErr)
		}
		return SignupOutput{}, err
	}

	return SignupOutput{RetryAt: reservation.RetryAt}, nil
}

func (s *Service) ResendSignup(ctx context.Context, email string) (SignupOutput, error) {
	emailHash := crypto.HashEmail(email)

	otp, err := crypto.GenerateOTP()
	if err != nil {
		return SignupOutput{}, fmt.Errorf("generate otp: %w", err)
	}

	otpHash := crypto.HashOTP(otp, s.cfg.OTPSecretKey)

	resend, err := s.challenges.Resend(ctx, emailHash, otpHash)
	if err != nil {
		return SignupOutput{}, err
	}
	if !resend.Reservation.Allowed {
		return SignupOutput{}, reservationError(resend.Reservation)
	}

	if err := s.dispatcher.DispatchSignupOTP(ctx, SignupOTP{
		Recipient: resend.Recipient,
		Email:     resend.Email,
		EmailHash: emailHash,
		OTP:       otp,
		OTPHash:   otpHash,
		ExpiresAt: resend.Reservation.ExpiresAt,
	}); err != nil {
		if releaseErr := s.challenges.Release(ctx, emailHash, otpHash); releaseErr != nil {
			return SignupOutput{}, errors.Join(err, releaseErr)
		}
		return SignupOutput{}, err
	}

	return SignupOutput{RetryAt: resend.Reservation.RetryAt}, nil
}

func (s *Service) VerifySignup(
	ctx context.Context,
	email, otp string,
) (domain.User, Tokens, error) {
	emailHash := crypto.HashEmail(email)
	challenge, err := s.challenges.Load(ctx, emailHash)
	if errors.Is(err, redis.Nil) {
		return domain.User{}, Tokens{}, domain.ErrChallengeExpired
	}
	if err != nil {
		return domain.User{}, Tokens{}, err
	}

	providedHash := crypto.HashOTP(otp, s.cfg.OTPSecretKey)
	if !hmac.Equal([]byte(challenge.OTPHash), []byte(providedHash)) {
		return domain.User{}, Tokens{}, s.recordFailedAttempt(ctx, emailHash)
	}

	user := domain.User{
		ID:           uuid.New(),
		Name:         challenge.Name,
		Email:        challenge.Email,
		PasswordHash: challenge.PasswordHash,
	}

	created, err := s.users.Create(ctx, user)
	if err != nil {
		return domain.User{}, Tokens{}, err
	}

	now := time.Now()
	accessExpiresAt := now.Add(s.cfg.AccessTokenExpiry)
	refreshExpiresAt := now.Add(s.cfg.RefreshTokenExpiry)

	accessToken, err := s.signer.SignAccess(user.ID.String(), now, accessExpiresAt)
	if err != nil {
		return domain.User{}, Tokens{}, err
	}
	refreshToken, err := s.signer.SignRefresh(user.ID.String(), now, refreshExpiresAt)
	if err != nil {
		return domain.User{}, Tokens{}, err
	}

	if err := s.challenges.Clear(ctx, emailHash); err != nil {
		return domain.User{}, Tokens{}, fmt.Errorf("clear signup challenge: %w", err)
	}

	return created, Tokens{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessExpiresAt,
		RefreshTokenExpiresAt: refreshExpiresAt,
	}, nil
}

func reservationError(r otpstore.Reservation) error {
	switch r.Reason {
	case otpstore.ReasonCooldownActive:
		return &domain.CooldownError{RetryAt: r.RetryAt}
	case otpstore.ReasonSendLimitReached:
		return &domain.SendLimitError{RetryAt: r.RetryAt}
	case otpstore.ReasonChallengeExpired:
		return domain.ErrChallengeExpired
	case otpstore.ReasonChallengeMismatch:
		return domain.ErrChallengeMismatch
	case otpstore.ReasonInvalidAttempts, otpstore.ReasonInvalidCooldown:
		return fmt.Errorf("identity: invalid reservation state: %s", r.Reason)
	default:
		return errors.New("identity: unexpected reservation state: " + r.Reason)
	}
}

func (s *Service) recordFailedAttempt(ctx context.Context, emailHash string) error {
	result, err := s.challenges.RecordFailedAttempt(ctx, emailHash)
	if err != nil {
		return err
	}

	switch result {
	case otpstore.AttemptExpired:
		return domain.ErrChallengeExpired
	case otpstore.AttemptRecorded:
		return domain.ErrOTPInvalid
	case otpstore.AttemptLimitReached:
		return domain.ErrTooManyAttempts
	default:
		return fmt.Errorf("identity: unexpected verification result: %d", result)
	}
}
