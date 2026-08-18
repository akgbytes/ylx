package handler

import (
	"crypto/hmac"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"

	"github.com/akgbytes/ylx/internal/auth"
	"github.com/akgbytes/ylx/internal/httpx"
	"github.com/akgbytes/ylx/internal/tasks"
)

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := zerolog.Ctx(ctx)

	var payload signupPayload
	if !httpx.DecodeJSON(w, r, &payload) {
		return
	}

	payload.normalize()

	if field, err := payload.validate(); err != nil {
		httpx.WriteValidationError(w, field, err.Error())
		return
	}

	// Check email availability
	query := `
		SELECT NOT EXISTS (
			SELECT 1
			FROM users
			WHERE email = $1
		)
	`

	var emailAvailable bool
	if err := h.db.QueryRowContext(ctx, query, payload.Email).Scan(&emailAvailable); err != nil {
		logger.Err(err).Msg("check email availability")
		httpx.WriteError(w, httpx.CodeInternal, "internal server error")
		return
	}

	if !emailAvailable {
		httpx.WriteError(w, httpx.CodeConflict, "email is already in use")
		return
	}

	passwordHash, err := auth.HashPassword(payload.Password, auth.DefaultArgonParams())
	if err != nil {
		logger.Err(err).Msg("hash password")
		httpx.WriteError(w, httpx.CodeInternal, "internal server error")
		return
	}

	otp, err := auth.GenerateOTP()
	if err != nil {
		logger.Err(err).Msg("generate OTP")
		httpx.WriteError(w, httpx.CodeInternal, "internal server error")
		return
	}

	otpHash := auth.HashOTP(otp, h.cfg.OTPSecretKey)

	challengeID, err := auth.GenerateRandomString(32)
	if err != nil {
		logger.Err(err).Msg("generate challenge ID")
		httpx.WriteError(w, httpx.CodeInternal, "internal server error")
		return
	}

	challenge := auth.SignupChallenge{
		Name:         payload.Name,
		Email:        payload.Email,
		PasswordHash: passwordHash,
		OTPHash:      otpHash,
	}

	// Reserving challenge in redis
	reservation, err := auth.ReserveSignupChallenge(
		ctx,
		h.rdb,
		payload.Email,
		challengeID,
		challenge,
		*h.cfg,
	)
	if err != nil {
		logger.Err(err).Msg("reserve signup challenge")
		httpx.WriteError(w, httpx.CodeInternal, "internal server error")
		return
	}
	if !reservation.Allowed {
		retryAt := reservation.RetryAt
		switch reservation.Reason {
		case auth.SignupReservationCooldownActive:
			httpx.WriteCooldownError(
				w,
				httpx.CodeOTPCooldownActive,
				"please wait before requesting another OTP",
				retryAt,
			)
		case auth.SignupReservationSendLimitReached:
			httpx.WriteRateLimitError(
				w,
				httpx.CodeOTPHourlyLimit,
				"too many verification code requests",
				h.cfg.OTPMaxSends, 0, retryAt,
			)
		default:
			logger.Error().Str("reason", reservation.Reason).Msg("invalid signup reservation state")
			httpx.WriteError(w, httpx.CodeInternal, "internal server error")
		}
		return
	}

	expiresAt := time.Now().Add(h.cfg.OTPExpiry)

	// Create email task and enqueue
	emailTask, err := tasks.NewSendSignUpOTPTask(tasks.SendSignUpOTPPayload{
		ChallengeID: challengeID,
		Email:       payload.Email,
		OTP:         otp,
		OTPHash:     otpHash,
		ExpiresAt:   expiresAt,
	})
	if err != nil {
		logger.Err(err).Msg("create signup OTP task")
		httpx.WriteError(w, httpx.CodeInternal, "internal server error")
		return
	}

	if _, err := h.asynqClient.EnqueueContext(
		ctx,
		emailTask,
		asynq.Queue("email"),
		asynq.MaxRetry(5),
		asynq.Timeout(15*time.Second),
		asynq.Deadline(expiresAt),
	); err != nil {
		// TODO: Release the signup reservation so the client can retry after an enqueue failure
		logger.Err(err).Msg("enqueue signup OTP task")
		httpx.WriteError(w, httpx.CodeInternal, "internal server error")
		return
	}

	httpx.WriteJSON(w, http.StatusAccepted, signUpResponse{
		ChallengeID: challengeID,
		RetryAt:     reservation.RetryAt,
	})
}

func (h *AuthHandler) ResendSignUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := zerolog.Ctx(ctx)

	var payload resendSignUpPayload
	if !httpx.DecodeJSON(w, r, &payload) {
		return
	}

	payload.normalize()

	if field, err := payload.validate(); err != nil {
		httpx.WriteValidationError(w, field, err.Error())
		return
	}

	otp, err := auth.GenerateOTP()
	if err != nil {
		logger.Err(err).Msg("generate signup OTP")
		httpx.WriteError(w, httpx.CodeInternal, "internal server error")
		return
	}

	otpHash := auth.HashOTP(otp, h.cfg.OTPSecretKey)
	reservation, err := auth.ResendSignupChallenge(ctx, h.rdb, payload.Email, payload.ChallengeID, otpHash, *h.cfg)
	if err != nil {
		logger.Err(err).Msg("resend signup challenge")
		httpx.WriteError(w, httpx.CodeInternal, "internal server error")
		return
	}

	if !reservation.Allowed {
		switch reservation.Reason {
		case auth.SignupReservationCooldownActive:
			httpx.WriteCooldownError(
				w,
				httpx.CodeOTPCooldownActive,
				"please wait before requesting another OTP",
				reservation.RetryAt,
			)
		case auth.SignupReservationSendLimitReached:
			httpx.WriteRateLimitError(
				w,
				httpx.CodeOTPHourlyLimit,
				"too many verification code requests",
				h.cfg.OTPMaxSends, 0, reservation.RetryAt,
			)
		case auth.SignupReservationChallengeExpired:
			httpx.WriteError(w, httpx.CodeOTPExpired, "verification code has expired")
		case auth.SignupReservationChallengeMismatch:
			httpx.WriteError(w, httpx.CodeOTPInvalid, "invalid verification code")
		default:
			logger.Error().Str("reason", reservation.Reason).Msg("invalid resend signup state")
			httpx.WriteError(w, httpx.CodeInternal, "internal server error")
		}
		return
	}

	expiresAt := time.Now().Add(h.cfg.OTPExpiry)

	// Create email task and enqueue
	emailTask, err := tasks.NewSendSignUpOTPTask(tasks.SendSignUpOTPPayload{
		ChallengeID: payload.ChallengeID,
		Email:       payload.Email,
		OTP:         otp,
		OTPHash:     otpHash,
		ExpiresAt:   expiresAt,
	})
	if err != nil {
		logger.Err(err).Msg("create signup OTP task")
		httpx.WriteError(w, httpx.CodeInternal, "internal server error")
		return
	}

	if _, err := h.asynqClient.EnqueueContext(
		ctx,
		emailTask,
		asynq.Queue("email"),
		asynq.MaxRetry(5),
		asynq.Timeout(15*time.Second),
		asynq.Deadline(expiresAt),
	); err != nil {
		// TODO: Release the signup reservation so the client can retry after an enqueue failure
		logger.Err(err).Msg("enqueue signup OTP task")
		httpx.WriteError(w, httpx.CodeInternal, "internal server error")
		return
	}

	httpx.WriteJSON(w, http.StatusAccepted, signUpResponse{
		ChallengeID: payload.ChallengeID,
		RetryAt:     reservation.RetryAt,
	})
}

func (h *AuthHandler) VerifySignUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := zerolog.Ctx(ctx)

	var payload verifySignUpPayload
	if !httpx.DecodeJSON(w, r, &payload) {
		return
	}

	payload.normalize()

	if field, err := payload.validate(); err != nil {
		httpx.WriteValidationError(w, field, err.Error())
		return
	}

	challenge, err := auth.LoadSignupChallenge(ctx, h.rdb, payload.ChallengeID)
	if errors.Is(err, redis.Nil) {
		httpx.WriteError(w, httpx.CodeOTPExpired, "verification code has expired")
		return
	}

	if err != nil {
		logger.Err(err).Msg("load signup challenge")
		httpx.WriteError(w, httpx.CodeInternal, "internal server error")
		return
	}

	if challenge.Email != payload.Email {
		httpx.WriteError(w, httpx.CodeOTPInvalid, "invalid verification code")
		return
	}

	providedOTPHash := auth.HashOTP(payload.OTP, h.cfg.OTPSecretKey)
	if !hmac.Equal([]byte(challenge.OTPHash), []byte(providedOTPHash)) {

		result, err := auth.RecordFailedSignupVerification(
			ctx,
			h.rdb,
			payload.ChallengeID,
			h.cfg.OTPMaxVerificationAttempts,
		)
		if err != nil {
			logger.Err(err).Msg("record failed signup verification")
			httpx.WriteError(w, httpx.CodeInternal, "internal server error")
			return
		}

		switch result {
		case auth.SignupVerificationExpired:
			httpx.WriteError(w, httpx.CodeOTPExpired, "verification code has expired")
		case auth.SignupVerificationRecorded:
			httpx.WriteError(w, httpx.CodeOTPInvalid, "invalid verification code")
		case auth.SignupVerificationLimitReached:
			httpx.WriteError(w, httpx.CodeTooManyRequests, "too many invalid verification attempts")
		default:
			logger.Error().Int64("result", int64(result)).Msg("invalid signup verification state")
			httpx.WriteError(w, httpx.CodeInternal, "internal server error")
		}

		return
	}

	userID := uuid.NewString()
	session, err := h.sessions.Issue(userID)
	if err != nil {
		logger.Err(err).Msg("issue session")
		httpx.WriteError(w, httpx.CodeInternal, "internal server error")
		return
	}

	// Create user
	query := `
		INSERT INTO users (id, name, email, password_hash, refresh_token_hash, refresh_token_expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, email, created_at
	`

	var user verifiedUserResponse
	err = h.db.QueryRowContext(
		ctx,
		query,
		userID,
		challenge.Name,
		challenge.Email,
		challenge.PasswordHash,
		session.RefreshTokenHash,
		session.RefreshTokenExpiresAt,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if cleanupErr := auth.DeleteSignupVerificationState(ctx, h.rdb, payload.ChallengeID); cleanupErr != nil {
				logger.Err(cleanupErr).Msg("delete signup verification state")
			}

			httpx.WriteError(w, httpx.CodeConflict, "email is already in use")
			return
		}

		logger.Err(err).Msg("create user")
		httpx.WriteError(w, httpx.CodeInternal, "internal server error")
		return
	}

	if err := auth.DeleteSignupVerificationState(ctx, h.rdb, payload.ChallengeID); err != nil {
		logger.Err(err).Msg("delete signup verification state")
	}

	logger.Info().Str("user_id", user.ID).Msg("user created")
	h.sessions.SetCookies(w, session)
	httpx.WriteJSON(w, http.StatusCreated, user)
}
