package domain

import (
	"time"
	"uuid"
)

type User struct {
	ID           uuid.UUID
	Name         string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
