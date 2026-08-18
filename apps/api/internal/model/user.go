package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                 uuid.UUID
	Name               string
	Email              string
	PasswordHash       string
	RefreshTokenHash   *string
	RefreshTokenExpiry *time.Time
	IsActive           bool
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
