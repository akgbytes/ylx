package rest

import (
	"time"

	"github.com/akgbytes/ylx/internal/listing/domain"
)

type listingResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       int64     `json:"price"`
	City        string    `json:"city"`
	SellerID    string    `json:"seller_id"`
	CategoryID  string    `json:"category_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func newListingResponses(listings []domain.Listing) []listingResponse {
	responses := make([]listingResponse, 0, len(listings))
	for _, listing := range listings {
		responses = append(responses, listingResponse{
			ID:          listing.ID.String(),
			Title:       listing.Title,
			Description: listing.Description,
			Price:       listing.Price,
			City:        listing.City,
			SellerID:    listing.SellerID.String(),
			CategoryID:  listing.CategoryID.String(),
			CreatedAt:   listing.CreatedAt,
			UpdatedAt:   listing.UpdatedAt,
		})
	}

	return responses
}
