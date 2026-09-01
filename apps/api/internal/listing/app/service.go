package app

import (
	"context"
	"fmt"

	"github.com/akgbytes/ylx/internal/listing/domain"
)

type Store interface {
	List(ctx context.Context) ([]domain.Listing, error)
}

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) List(ctx context.Context) ([]domain.Listing, error) {
	listings, err := s.store.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list listings: %w", err)
	}

	return listings, nil
}
