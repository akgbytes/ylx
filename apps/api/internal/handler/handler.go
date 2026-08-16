package handler

import (
	"database/sql"
)

type Handlers struct {
	Health  *HealthHandler
	Listing *ListingHandler
	Auth    *AuthHandler
}

func NewHandlers(db *sql.DB) *Handlers {
	return &Handlers{
		Health:  NewHealthHandler(),
		Listing: NewListingHandler(db),
		Auth:    NewAuthHandler(),
	}
}
