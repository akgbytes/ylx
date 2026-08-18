package handler

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/rs/zerolog"

	"github.com/akgbytes/ylx/internal/auth"
	"github.com/akgbytes/ylx/internal/httpx"
	"github.com/akgbytes/ylx/internal/model"
)

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := zerolog.Ctx(ctx)

	var payload signInPayload
	if !httpx.DecodeJSON(w, r, &payload) {
		return
	}

	payload.normalize()

	if field, err := payload.validate(); err != nil {
		httpx.WriteValidationError(w, field, err.Error())
		return
	}

	query := `
		SELECT id, name, email, password_hash, is_active
		FROM users
		WHERE email = $1
	`

	var user model.User
	err := h.db.QueryRowContext(ctx, query, payload.Email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.IsActive,
	)

	if errors.Is(err, sql.ErrNoRows) {
		httpx.WriteError(w, httpx.CodeUnauthorized, "invalid email or password")
		return
	}
	if err != nil {
		logger.Err(err).Msg("find user for sign in")
		httpx.WriteError(w, httpx.CodeInternal, "internal server error")
		return
	}

	passwordMatches, err := auth.VerifyPassword(payload.Password, user.PasswordHash)
	if err != nil {
		logger.Err(err).Str("user_id", user.ID.String()).Msg("verify user password")
		httpx.WriteError(w, httpx.CodeInternal, "internal server error")
		return
	}
	if !passwordMatches {
		httpx.WriteError(w, httpx.CodeUnauthorized, "invalid email or password")
		return
	}

	if !user.IsActive {
		httpx.WriteError(w, httpx.CodeForbidden, "account is inactive")
		return
	}

	session, err := h.sessions.Issue(user.ID.String())
	if err != nil {
		logger.Err(err).Str("user_id", user.ID.String()).Msg("issue session")
		httpx.WriteError(w, httpx.CodeInternal, "internal server error")
		return
	}

	query = `
		UPDATE users
		SET refresh_token_hash = $1, refresh_token_expires_at = $2
		WHERE id = $3 AND is_active = TRUE
	`

	result, err := h.db.ExecContext(ctx, query, session.RefreshTokenHash, session.RefreshTokenExpiresAt, user.ID)
	if err != nil {
		logger.Err(err).Str("user_id", user.ID.String()).Msg("update refresh token")
		httpx.WriteError(w, httpx.CodeInternal, "internal server error")
		return

	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {

		logger.Err(err).Str("user_id", user.ID.String()).Msg("read refresh token update count")
		httpx.WriteError(w, httpx.CodeInternal, "internal server error")
		return
	}

	if rowsAffected != 1 {
		logger.Error().
			Str("user_id", user.ID.String()).
			Int64("rows_affected", rowsAffected).
			Msg("update refresh token affected unexpected row count")
		httpx.WriteError(w, httpx.CodeInternal, "internal server error")
		return
	}

	logger.Info().Str("user_id", user.ID.String()).Msg("user signed in")
	h.sessions.SetCookies(w, session)
	httpx.WriteJSON(w, http.StatusOK, signedInUserResponse{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := zerolog.Ctx(ctx)

	h.sessions.ClearCookies(w)

	accessCookie, accessCookieErr := r.Cookie(h.cfg.AccessTokenName)
	if accessCookieErr == nil {
		userID, err := h.sessions.VerifyAccessToken(accessCookie.Value)
		if err == nil {
			query := `
				UPDATE users
				SET refresh_token_hash = NULL, refresh_token_expires_at = NULL
				WHERE id = $1 AND is_active = TRUE
			`

			if _, err := h.db.ExecContext(ctx, query, userID); err != nil {
				logger.Err(err).Str("user_id", userID).Msg("revoke user session")
				httpx.WriteError(w, httpx.CodeInternal, "internal server error")
				return
			}

			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	refreshCookie, refreshCookieErr := r.Cookie(h.cfg.RefreshTokenName)
	if refreshCookieErr == nil {
		query := `
			UPDATE users
			SET refresh_token_hash = NULL, refresh_token_expires_at = NULL
			WHERE refresh_token_hash = $1
			  AND refresh_token_expires_at > CURRENT_TIMESTAMP
		`

		if _, err := h.db.ExecContext(ctx, query, auth.HashRefreshToken(refreshCookie.Value)); err != nil {
			logger.Err(err).Msg("revoke refresh-token session")
			httpx.WriteError(w, httpx.CodeInternal, "internal server error")
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := zerolog.Ctx(ctx)

	refreshCookie, err := r.Cookie(h.cfg.RefreshTokenName)
	if err != nil {
		h.sessions.ClearCookies(w)
		httpx.WriteError(w, httpx.CodeUnauthorized, "invalid or expired session")
		return
	}

	refreshTokenHash := auth.HashRefreshToken(refreshCookie.Value)
	query := `
		SELECT id
		FROM users
		WHERE refresh_token_hash = $1
		  AND refresh_token_expires_at > CURRENT_TIMESTAMP
		  AND is_active = TRUE
	`

	var userID string
	err = h.db.QueryRowContext(ctx, query, refreshTokenHash).Scan(&userID)
	if errors.Is(err, sql.ErrNoRows) {
		h.sessions.ClearCookies(w)
		httpx.WriteError(w, httpx.CodeUnauthorized, "invalid or expired session")
		return
	}
	if err != nil {
		logger.Err(err).Msg("find user by refresh token")
		httpx.WriteError(w, httpx.CodeInternal, "internal server error")
		return
	}

	session, err := h.sessions.Issue(userID)
	if err != nil {
		logger.Err(err).Str("user_id", userID).Msg("issue session")
		httpx.WriteError(w, httpx.CodeInternal, "internal server error")
		return
	}

	query = `
		UPDATE users
		SET refresh_token_hash = $1,
			refresh_token_expires_at = $2
		WHERE id = $3
		  AND refresh_token_hash = $4
		  AND refresh_token_expires_at > CURRENT_TIMESTAMP
		  AND is_active = TRUE
	`

	result, err := h.db.ExecContext(
		ctx,
		query,
		session.RefreshTokenHash,
		session.RefreshTokenExpiresAt,
		userID,
		refreshTokenHash,
	)
	if err != nil {
		logger.Err(err).Str("user_id", userID).Msg("rotate refresh token")
		httpx.WriteError(w, httpx.CodeInternal, "internal server error")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Err(err).Str("user_id", userID).Msg("read refresh token update count")
		httpx.WriteError(w, httpx.CodeInternal, "internal server error")
		return
	}
	if rowsAffected != 1 {
		h.sessions.ClearCookies(w)
		httpx.WriteError(w, httpx.CodeUnauthorized, "invalid or expired session")
		return
	}

	h.sessions.SetCookies(w, session)
	w.WriteHeader(http.StatusNoContent)
}
