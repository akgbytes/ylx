package model

import "time"

type Listing struct {
	ID          string
	Title       string
	Description string
	Price       int64
	City        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
