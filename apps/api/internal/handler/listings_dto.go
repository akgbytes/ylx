package handler

import "github.com/akgbytes/ylx/internal/model"

type ListListingsResponse struct {
	Listings []model.Listing `json:"listings"`
}

type ListingNotFoundResponse struct {
	Message string `json:"message"`
}
