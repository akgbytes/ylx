package handler

import "github.com/akgbytes/ylx/internal/model"

type ListingsResponse struct {
	Listings []model.Listing `json:"listings"`
}

type ListingsErrorResponse struct {
	Message string `json:"message"`
}
