package handler

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/akgbytes/ylx/internal/httpx"
)

type ListingHandler struct {
	db *sql.DB
}

func NewListingHandler(db *sql.DB) *ListingHandler {
	return &ListingHandler{
		db: db,
	}
}

func (h *ListingHandler) Create(w http.ResponseWriter, r *http.Request) {
	logger := zerolog.Ctx(r.Context())

	var p createListingPayload
	if !httpx.DecodeJSON(w, r, &p) {
		return
	}

	p.normalize()

	if field, err := p.validate(); err != nil {
		httpx.WriteValidationError(w, field, err.Error())
		return
	}

	query := `
		INSERT INTO listings (title, description, price, city)
		VALUES ($1, $2, $3, $4)
		RETURNING id, title, description, price, city, created_at, updated_at
	`

	var listing listingResponse
	err := h.db.QueryRowContext(
		r.Context(),
		query,
		p.Title, p.Description, p.Price, p.City,
	).Scan(
		&listing.ID,
		&listing.Title,
		&listing.Description,
		&listing.Price,
		&listing.City,
		&listing.CreatedAt,
		&listing.UpdatedAt,
	)
	if err != nil {
		logger.Err(err).Msg("create listing")
		httpx.WriteError(w, httpx.CodeInternal, "internal server error")
		return
	}

	logger.Info().Str("listing_id", listing.ID).Msg("listing created")
	httpx.WriteJSON(w, http.StatusCreated, listing)
}

func (h *ListingHandler) List(w http.ResponseWriter, r *http.Request) {
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
			logger.Err(err).Msg("close listings rows")
		}
	}()

	listings := []listingResponse{}
	for rows.Next() {
		var listing listingResponse
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
		logger.Err(err).Msg("iterate listings rows")
		httpx.WriteError(w, httpx.CodeInternal, "internal server error")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, listings)
}

func (h *ListingHandler) Delete(w http.ResponseWriter, r *http.Request) {
	logger := zerolog.Ctx(r.Context())
	listingID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		httpx.WriteError(w, httpx.CodeBadRequest, "invalid listing id")
		return
	}

	query := `
		DELETE FROM listings
		WHERE id = $1;
	`

	result, err := h.db.ExecContext(r.Context(), query, listingID)
	if err != nil {
		logger.Err(err).Any("listing_id", listingID).Msg("delete listing")
		httpx.WriteError(w, httpx.CodeInternal, "internal server error")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Err(err).Any("listing_id", listingID).Msg("get deleted listing count")
		httpx.WriteError(w, httpx.CodeInternal, "internal server error")
		return
	}

	if rowsAffected == 0 {
		httpx.WriteError(w, httpx.CodeNotFound, "listing not found")
		return
	}

	logger.Info().Any("listing_id", listingID).Msg("listing deleted")

	w.WriteHeader(http.StatusNoContent)
}
