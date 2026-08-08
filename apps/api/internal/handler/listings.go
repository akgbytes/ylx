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

func (h *ListingsHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	query := `
		SELECT id, title, description, price, city, created_at, updated_at
		FROM listings
		ORDER BY title DESC;
	`

	rows, err := h.db.QueryContext(ctx, query)
	if err != nil {
		h.logger.ErrorContext(ctx, "failed to query listings", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := rows.Close(); err != nil {
			h.logger.ErrorContext(ctx, "failed to close listing rows", "error", err)
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
			&listing.Created_At,
			&listing.Updated_At,
		); err != nil {
			h.logger.ErrorContext(ctx, "failed to scan listing", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		listings = append(listings, listing)
	}

	if err := rows.Err(); err != nil {
		h.logger.ErrorContext(ctx, "failed while iterating listings", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(ListingsResponse{
		Listings: listings,
	})
}

func (h *ListingsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")

	query := `
		DELETE FROM listings
		WHERE id = $1;
	`

	result, err := h.db.ExecContext(ctx, query, id)
	if err != nil {
		h.logger.ErrorContext(ctx,
			"failed to delete listing",
			"listing_id", id,
			"error", err,
		)

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		h.logger.ErrorContext(ctx,
			"failed to get deleted listing count",
			"listing_id", id,
			"error", err,
		)

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		_ = json.NewEncoder(w).Encode(ListingsErrorResponse{
			Message: "listing not found",
		})

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
