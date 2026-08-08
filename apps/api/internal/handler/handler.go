package handler

import (
	"database/sql"
	"log/slog"
)

type Handlers struct {
	Health  *HealthHandler
	Listing *ListingsHandler
}

func New(db *sql.DB, logger *slog.Logger) *Handlers {
	return &Handlers{
		Health:  NewHealthHandler(),
		Listing: NewListingsHandler(db, logger),
	}
}
