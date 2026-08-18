package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/akgbytes/ylx/internal/config"
)

type SessionManager struct {
	accessTokenName    string
	refreshTokenName   string
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
	jwtSecretKey       []byte
	secure             bool
}

type Session struct {
	AccessToken           string
	RefreshToken          string
	RefreshTokenHash      string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
}

func NewSessionManager(cfg *config.AuthConfig, secure bool) *SessionManager {
	return &SessionManager{
		accessTokenName:    cfg.AccessTokenName,
		refreshTokenName:   cfg.RefreshTokenName,
		accessTokenExpiry:  cfg.AccessTokenExpiry,
		refreshTokenExpiry: cfg.RefreshTokenExpiry,
		jwtSecretKey:       []byte(cfg.JWTSecretKey),
		secure:             secure,
	}
}

func (s *SessionManager) Issue(userID string) (Session, error) {
	now := time.Now()
	accessTokenExpiresAt := now.Add(s.accessTokenExpiry)
	refreshTokenExpiresAt := now.Add(s.refreshTokenExpiry)

	accessToken, err := s.signAccessToken(userID, now, accessTokenExpiresAt)
	if err != nil {
		return Session{}, err
	}

	refreshToken, err := newRefreshToken()
	if err != nil {
		return Session{}, err
	}

	return Session{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		RefreshTokenHash:      HashRefreshToken(refreshToken),
		AccessTokenExpiresAt:  accessTokenExpiresAt,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
	}, nil
}

func (m *SessionManager) VerifyAccessToken(accessToken string) (string, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(
		accessToken,
		claims,
		func(token *jwt.Token) (any, error) {
			return m.jwtSecretKey, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer("ylx"),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
	)
	if err != nil {
		return "", fmt.Errorf("parse access token: %w", err)
	}

	if !token.Valid {
		return "", errors.New("access token is invalid")
	}

	if _, err := uuid.Parse(claims.Subject); err != nil {
		return "", fmt.Errorf("parse access token subject: %w", err)
	}

	return claims.Subject, nil
}

//nolint:gosec // Cookie security is enabled in prod  configurable for local development.
func (m *SessionManager) SetCookies(w http.ResponseWriter, session Session) {
	http.SetCookie(w, &http.Cookie{
		Name:     m.accessTokenName,
		Value:    session.AccessToken,
		Path:     "/",
		Expires:  session.AccessTokenExpiresAt,
		MaxAge:   int(m.accessTokenExpiry.Seconds()),
		HttpOnly: true,
		Secure:   m.secure,
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     m.refreshTokenName,
		Value:    session.RefreshToken,
		Path:     "/",
		Expires:  session.RefreshTokenExpiresAt,
		MaxAge:   int(m.refreshTokenExpiry.Seconds()),
		HttpOnly: true,
		Secure:   m.secure,
		SameSite: http.SameSiteLaxMode,
	})
}

//nolint:gosec // Cookie security is enabled in prod  configurable for local development.
func (m *SessionManager) ClearCookies(w http.ResponseWriter) {
	expiresAt := time.Unix(1, 0)

	http.SetCookie(w, &http.Cookie{
		Name:     m.accessTokenName,
		Value:    "",
		Path:     "/",
		Expires:  expiresAt,
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   m.secure,
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     m.refreshTokenName,
		Value:    "",
		Path:     "/",
		Expires:  expiresAt,
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   m.secure,
		SameSite: http.SameSiteLaxMode,
	})
}

func (m *SessionManager) signAccessToken(userID string, issuedAt, expiresAt time.Time) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "ylx",
		Subject:   userID,
		IssuedAt:  jwt.NewNumericDate(issuedAt),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(m.jwtSecretKey)
	if err != nil {
		return "", fmt.Errorf("sign access token: %w", err)
	}

	return signedToken, nil
}

func newRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate refresh token: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func HashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
