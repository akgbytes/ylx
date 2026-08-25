package domain

import (
	"errors"
	"time"
)

var (
	ErrEmailTaken         = errors.New("identity: email already in use")
	ErrInvalidCredentials = errors.New("identity: invalid email or password")
	ErrUserNotFound       = errors.New("identity: user not found")
	ErrChallengeExpired   = errors.New("identity: verification code has expired")
	ErrChallengeMismatch  = errors.New("identity: verification code is invalid")
	ErrOTPInvalid         = errors.New("identity: verification code is invalid")
	ErrTooManyAttempts    = errors.New("identity: too many invalid verification attempts")
)

type CooldownError struct {
	RetryAt time.Time
}

func (e *CooldownError) Error() string {
	return "identity: verification code cooldown is active"
}

type SendLimitError struct {
	RetryAt time.Time
}

func (e *SendLimitError) Error() string {
	return "identity: verification code send limit reached"
}
