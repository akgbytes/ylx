package handler

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/akgbytes/ylx/internal/model"
)

type ListingsHandler struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewListingsHandler(db *sql.DB, logger *slog.Logger) *ListingsHandler {
	return &ListingsHandler{
		db:     db,
		logger: logger,
	}
}

func (handler *ListingsHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	query := `
		SELECT id, title, description, price, city, created_at, updated_at
		FROM listings
		ORDER BY title DESC;
	`

	rows, err := handler.db.QueryContext(ctx, query)
	if err != nil {
		handler.logger.Error("query listings", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := rows.Close(); err != nil {
			handler.logger.Error("close listing rows", "error", err)
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
			handler.logger.Error("scan listing", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		listings = append(listings, listing)
	}

	if err := rows.Err(); err != nil {
		handler.logger.Error("iterate listing rows", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(ListListingsResponse{
		Listings: listings,
	}); err != nil {
		handler.logger.Error("write listings response", "error", err)
	}
}

func (handler *ListingsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	listingID := r.PathValue("id")

	query := `
		DELETE FROM listings
		WHERE id = $1;
	`

	result, err := handler.db.ExecContext(ctx, query, listingID)
	if err != nil {
		handler.logger.Error(
			"delete listing",
			"listing_id", listingID,
			"error", err,
		)

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		handler.logger.Error(
			"get deleted listing count",
			"listing_id", listingID,
			"error", err,
		)

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

	w.WriteHeader(http.StatusNoContent)
}
