package domain

import "errors"

var (
	ErrEmailTaken         = errors.New("identity: email already in use")
	ErrInvalidCredentials = errors.New("identity: invalid email or password")
	ErrUserNotFound       = errors.New("identity: user not found")
)
