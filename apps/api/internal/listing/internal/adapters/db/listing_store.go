package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/akgbytes/ylx/internal/listing/internal/domain"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) List(ctx context.Context) ([]domain.Listing, error) {
	query := `
		SELECT l.id, l.seller_id, l.category_id, c.name, l.title, l.description, l.price, l.city, l.created_at, l.updated_at
		FROM listings l
		JOIN categories c ON c.id = l.category_id
		ORDER BY l.created_at DESC, l.id DESC
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query listings: %w", err)
	}

	defer func() { _ = rows.Close() }()

	listings := []domain.Listing{}

	for rows.Next() {
		listing, err := scanListing(rows)
		if err != nil {
			return nil, fmt.Errorf("scan listing: %w", err)
		}

		listings = append(listings, listing)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate listings: %w", err)
	}
	return listings, nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanListing(row rowScanner) (domain.Listing, error) {
	var listing domain.Listing
	err := row.Scan(
		&listing.ID,
		&listing.SellerID,
		&listing.CategoryID,
		&listing.CategoryName,
		&listing.Title,
		&listing.Description,
		&listing.Price,
		&listing.City,
		&listing.CreatedAt,
		&listing.UpdatedAt,
	)
	return listing, err
}
