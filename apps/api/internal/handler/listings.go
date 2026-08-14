package handler

import (
	"database/sql"
	"net/http"

	"github.com/rs/zerolog"

	"github.com/akgbytes/ylx/internal/httpx"
	"github.com/akgbytes/ylx/internal/model"
)

type ListingsHandler struct {
	db *sql.DB
}

func NewListingsHandler(db *sql.DB) *ListingsHandler {
	return &ListingsHandler{
		db: db,
	}
}

func (h *ListingsHandler) List(w http.ResponseWriter, r *http.Request) {
	logger := zerolog.Ctx(r.Context())
	query := `
		SELECT id, title, description, price, city, created_at, updated_at
		FROM listings
		ORDER BY title DESC;
	`

	rows, err := h.db.QueryContext(r.Context(), query)
	if err != nil {
		logger.Err(err).Msg("query listings")
		httpx.WriteError(w, httpx.CodeInternal, "internal server error")
		return
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Err(err).Msg("close listing rows")
		}
	}()

	listings := []model.Listing{}

	for rows.Next() {
		var listing model.Listing

		if err := rows.Scan(
			&listing.ID,
			&listing.Title,
			&listing.Description,
			&listing.Price,
			&listing.City,
			&listing.CreatedAt,
			&listing.UpdatedAt,
		); err != nil {
			logger.Err(err).Msg("scan listing")
			httpx.WriteError(w, httpx.CodeInternal, "internal server error")
			return
		}
		listings = append(listings, listing)
	}

	if err := rows.Err(); err != nil {
		logger.Err(err).Msg("iterate listing rows")
		httpx.WriteError(w, httpx.CodeInternal, "internal server error")
		return
	}

	logger.Info().Msg("listing fetched")

	httpx.WriteJSON(w, http.StatusOK, listings, nil)
}

func (h *ListingsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	logger := zerolog.Ctx(r.Context())
	listingID := r.PathValue("id")

	query := `
		DELETE FROM listings
		WHERE id = $1;
	`

	result, err := h.db.ExecContext(r.Context(), query, listingID)
	if err != nil {
		logger.Err(err).Str("listing_id", listingID).Msg("delete listing")
		httpx.WriteError(w, httpx.CodeInternal, "internal server error")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Err(err).Str("listing_id", listingID).Msg("get deleted listing count")
		httpx.WriteError(w, httpx.CodeInternal, "internal server error")
		return
	}

	if rowsAffected == 0 {
		httpx.WriteError(w, httpx.CodeNotFound, "listing not found")
		return
	}

	logger.Info().Str("listing_id", listingID).Msg("listing deleted")

	w.WriteHeader(http.StatusNoContent)
}
