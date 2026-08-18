package model

import (
	"time"

	"github.com/google/uuid"
)

type Listing struct {
	ID          uuid.UUID
	Title       string
	Description string
	Price       int64
	City        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
