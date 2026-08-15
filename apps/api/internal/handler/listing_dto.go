package handler

import (
	"errors"
	"strings"
	"time"
)

type createListingPayload struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	City        string `json:"city"`
}

type listingResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       int64     `json:"price"`
	City        string    `json:"city"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (p *createListingPayload) validate() (string, error) {
	if strings.TrimSpace(p.Title) == "" {
		return "title", errors.New("title is required")
	}

	if strings.TrimSpace(p.Description) == "" {
		return "description", errors.New("description is required")
	}

	if strings.TrimSpace(p.City) == "" {
		return "city", errors.New("city is required")
	}

	if p.Price <= 0 {
		return "price", errors.New("price must be greater than 0")
	}

	return "", nil
}
