package handler

import (
	"database/sql"
	"log/slog"
)

type Handlers struct {
	Health   *HealthHandler
	Listings *ListingsHandler
}

func NewHandlers(db *sql.DB, logger *slog.Logger) *Handlers {
	return &Handlers{
		Health:   NewHealthHandler(),
		Listings: NewListingsHandler(db, logger),
	}
}
