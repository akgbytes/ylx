package app

import (
	"context"
	"fmt"
	"time"
	"uuid"

	"github.com/akgbytes/ylx/internal/identity/adapters/crypto"
	"github.com/akgbytes/ylx/internal/identity/domain"
)

type Tokens struct {
	AccessToken           string
	RefreshToken          string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
}

func (s *Service) SignIn(ctx context.Context, email, password string) (domain.User, Tokens, error) {
	user, err := s.users.ByEmail(ctx, email)
	if err != nil {
		return domain.User{}, Tokens{}, err
	}

	matches, err := crypto.VerifyPassword(password, user.PasswordHash)
	if err != nil {
		return domain.User{}, Tokens{}, fmt.Errorf("verify password: %w", err)
	}

	if !matches {
		return domain.User{}, Tokens{}, domain.ErrInvalidCredentials
	}

	now := time.Now()
	accessExpiresAt := now.Add(s.cfg.AccessTokenExpiry)
	refreshExpiresAt := now.Add(s.cfg.RefreshTokenExpiry)

	accessToken, err := s.signer.SignAccess(user.ID.String(), now, accessExpiresAt)
	if err != nil {
		return domain.User{}, Tokens{}, err
	}
	refreshToken, err := s.signer.SignRefresh(user.ID.String(), now, refreshExpiresAt)
	if err != nil {
		return domain.User{}, Tokens{}, err
	}

	return user, Tokens{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessExpiresAt,
		RefreshTokenExpiresAt: refreshExpiresAt,
	}, nil
}

func (s *Service) Refresh(ctx context.Context, userID string) (Tokens, error) {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return Tokens{}, fmt.Errorf("parse user id: %w", err)
	}

	if _, err := s.users.ByID(ctx, parsedUserID); err != nil {
		return Tokens{}, err
	}

	now := time.Now()
	accessExpiresAt := now.Add(s.cfg.AccessTokenExpiry)
	refreshExpiresAt := now.Add(s.cfg.RefreshTokenExpiry)

	accessToken, err := s.signer.SignAccess(userID, now, accessExpiresAt)
	if err != nil {
		return Tokens{}, err
	}

	refreshToken, err := s.signer.SignRefresh(userID, now, refreshExpiresAt)
	if err != nil {
		return Tokens{}, err
	}

	return Tokens{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessExpiresAt,
		RefreshTokenExpiresAt: refreshExpiresAt,
	}, nil
}
