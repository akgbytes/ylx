package token

import (
	"errors"
	"fmt"
	"time"
	"uuid"

	"github.com/golang-jwt/jwt/v5"
)

const issuer = "ylx"

type Claims struct {
	UserID string
}

type Signer struct {
	accessKey  []byte
	refreshKey []byte
}

func NewSigner(accessKey, refreshKey string) *Signer {
	return &Signer{
		accessKey:  []byte(accessKey),
		refreshKey: []byte(refreshKey),
	}
}

func (s *Signer) SignAccess(userID string, issuedAt, expiresAt time.Time) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    issuer,
		Subject:   userID,
		IssuedAt:  jwt.NewNumericDate(issuedAt),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	}

	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.accessKey)
	if err != nil {
		return "", fmt.Errorf("sign access token: %w", err)
	}

	return signed, nil
}

func (s *Signer) VerifyAccess(accessToken string) (Claims, error) {
	claims := &jwt.RegisteredClaims{}

	parsed, err := jwt.ParseWithClaims(
		accessToken,
		claims,
		func(t *jwt.Token) (any, error) { return s.accessKey, nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(issuer),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
	)
	if err != nil {
		return Claims{}, fmt.Errorf("parse access token: %w", err)
	}

	if !parsed.Valid {
		return Claims{}, errors.New("access token is invalid")
	}

	if _, err := uuid.Parse(claims.Subject); err != nil {
		return Claims{}, fmt.Errorf("parse access token subject: %w", err)
	}

	return Claims{UserID: claims.Subject}, nil
}

func (s *Signer) SignRefresh(userID string, issuedAt, expiresAt time.Time) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    issuer,
		Subject:   userID,
		IssuedAt:  jwt.NewNumericDate(issuedAt),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	}

	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.refreshKey)
	if err != nil {
		return "", fmt.Errorf("sign refresh token: %w", err)
	}

	return signed, nil
}

func (s *Signer) VerifyRefresh(refreshToken string) (Claims, error) {
	claims := &jwt.RegisteredClaims{}

	parsed, err := jwt.ParseWithClaims(
		refreshToken,
		claims,
		func(t *jwt.Token) (any, error) { return s.refreshKey, nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(issuer),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
	)
	if err != nil {
		return Claims{}, fmt.Errorf("parse refresh token: %w", err)
	}

	if !parsed.Valid {
		return Claims{}, errors.New("refresh token is invalid")
	}

	if _, err := uuid.Parse(claims.Subject); err != nil {
		return Claims{}, fmt.Errorf("parse access token subject: %w", err)
	}

	return Claims{UserID: claims.Subject}, nil
}
