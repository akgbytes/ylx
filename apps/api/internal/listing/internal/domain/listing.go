package domain

import (
	"time"

	"github.com/google/uuid"
)

type Listing struct {
	ID           uuid.UUID
	Title        string
	Description  string
	City         string
	Price        int64
	SellerID     uuid.UUID
	CategoryID   uuid.UUID
	CategoryName string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
