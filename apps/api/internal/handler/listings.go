package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog"

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
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Err(err).Msg("close listing rows")
		}
	}()

	var listings []model.Listing

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
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		listings = append(listings, listing)
	}

	if err := rows.Err(); err != nil {
		logger.Err(err).Msg("iterate listing rows")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	logger.Info().Msg("listing fetched")

	_ = json.NewEncoder(w).Encode(ListListingsResponse{
		Listings: listings,
	})
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
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Err(err).Str("listing_id", listingID).Msg("get deleted listing count")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		_ = json.NewEncoder(w).Encode(ListingNotFoundResponse{
			Message: "listing not found",
		})

		return
	}

	logger.Info().Str("listing_id", listingID).Msg("listing deleted")

	w.WriteHeader(http.StatusNoContent)
}
