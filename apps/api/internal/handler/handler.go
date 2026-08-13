package handler

import (
	"database/sql"
)

type Handlers struct {
	Health   *HealthHandler
	Listings *ListingsHandler
}

func NewHandlers(db *sql.DB) *Handlers {
	return &Handlers{
		Health:   NewHealthHandler(),
		Listings: NewListingsHandler(db),
	}
}
