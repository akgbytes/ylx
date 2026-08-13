package handler

import (
	"database/sql"

	"github.com/rs/zerolog"
)

type Handlers struct {
	Health   *HealthHandler
	Listings *ListingsHandler
}

func NewHandlers(db *sql.DB, logger zerolog.Logger) *Handlers {
	return &Handlers{
		Health:   NewHealthHandler(),
		Listings: NewListingsHandler(db, logger),
	}
}
