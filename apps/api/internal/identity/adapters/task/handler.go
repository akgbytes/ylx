package task

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"

	"github.com/akgbytes/ylx/internal/identity/adapters/otpstore"
	"github.com/akgbytes/ylx/internal/platform/mailer"
)

type Handler struct {
	sender     mailer.Sender
	challenges *otpstore.Store
	logger     zerolog.Logger
}

func NewHandler(sender mailer.Sender, challenges *otpstore.Store, logger zerolog.Logger) *Handler {
	return &Handler{
		sender:     sender,
		challenges: challenges,
		logger:     logger,
	}
}

func (h *Handler) Register(mux *asynq.ServeMux) {
	mux.HandleFunc(TypeSignupOTP, h.handleSignupOTP)
}

func (h *Handler) handleSignupOTP(ctx context.Context, task *asynq.Task) error {
	var payload signupOTPTaskPayload

	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("decode signup otp task: %w", asynq.SkipRetry)
	}

	if payload.Email == "" || payload.EmailHash == "" ||
		payload.OTP == "" || payload.OTPHash == "" ||
		payload.Recipient == "" || payload.ExpiresAt.IsZero() {
		return fmt.Errorf("validate signup otp task: %w", asynq.SkipRetry)
	}

	if !time.Now().Before(payload.ExpiresAt) {
		return fmt.Errorf("signup otp expired: %w", asynq.SkipRetry)
	}

	if err := h.ensureCurrent(ctx, payload.EmailHash, payload.OTPHash); err != nil {
		return err
	}

	h.logger.Info().
		Str("email", payload.Email).
		Msg("sending signup otp email")

	if err := h.sender.Send(ctx, payload.Email, signupOTPTemplate(payload.OTP, payload.Recipient)); err != nil {
		return fmt.Errorf("send signup otp: %w", err)
	}

	h.logger.Info().
		Str("email", payload.Email).
		Msg("signup otp email sent")

	return nil
}

func (h *Handler) ensureCurrent(ctx context.Context, emailHash, otpHash string) error {
	challenge, err := h.challenges.Load(ctx, emailHash)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return fmt.Errorf("signup challenge expired: %w", asynq.SkipRetry)
		}
		return fmt.Errorf("load signup challenge: %w", err)
	}

	if challenge.OTPHash != otpHash {
		return fmt.Errorf("signup otp superseded: %w", asynq.SkipRetry)
	}

	return nil
}
