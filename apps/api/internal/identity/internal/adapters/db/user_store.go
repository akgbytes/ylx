package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"uuid"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/akgbytes/ylx/internal/identity/internal/domain"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) EmailTaken(ctx context.Context, email string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM users
			WHERE email = $1
		)
	`

	var taken bool
	if err := s.db.QueryRowContext(ctx, query, email).Scan(&taken); err != nil {
		return false, fmt.Errorf("check email availability: %w", err)
	}

	return taken, nil
}

func (s *UserStore) Create(ctx context.Context, user domain.User) (domain.User, error) {
	query := `
		INSERT INTO users (id, name, email, password_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, email, password_hash, created_at, updated_at
	`

	created, err := scanUser(s.db.QueryRowContext(ctx, query, user.ID, user.Name, user.Email, user.PasswordHash))
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.User{}, domain.ErrEmailTaken
		}

		return domain.User{}, fmt.Errorf("create user: %w", err)
	}

	return created, nil
}

func (s *UserStore) ByEmail(ctx context.Context, email string) (domain.User, error) {
	query := `
		SELECT id, name, email, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user, err := scanUser(s.db.QueryRowContext(ctx, query, email))
	if errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, domain.ErrUserNotFound
	}

	if err != nil {
		return domain.User{}, fmt.Errorf("find user by email: %w", err)
	}

	return user, nil
}

func (s *UserStore) ByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	query := `
		SELECT id, name, email, password_hash, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user, err := scanUser(s.db.QueryRowContext(ctx, query, id))
	if errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, domain.ErrUserNotFound
	}

	if err != nil {
		return domain.User{}, fmt.Errorf("find user by id: %w", err)
	}

	return user, nil
}

func scanUser(row *sql.Row) (domain.User, error) {
	var user domain.User
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	return user, err
}
