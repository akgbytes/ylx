package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"

	"github.com/akgbytes/ylx/internal/auth"
	"github.com/akgbytes/ylx/internal/email"
	"github.com/akgbytes/ylx/internal/tasks"
)

type EmailWorker struct {
	sender email.Sender
	rdb    *redis.Client
}

func NewEmailWorker(sender email.Sender, rdb *redis.Client) *EmailWorker {
	return &EmailWorker{sender: sender, rdb: rdb}
}

func (w *EmailWorker) HandleSendSignUpOTP(ctx context.Context, asynqTask *asynq.Task) error {
	var payload tasks.SendSignUpOTPPayload

	if err := json.Unmarshal(asynqTask.Payload(), &payload); err != nil {
		return fmt.Errorf("decode signup OTP task: %w", asynq.SkipRetry)
	}

	if payload.ChallengeID == "" ||
		payload.Email == "" ||
		payload.OTP == "" ||
		payload.OTPHash == "" ||
		payload.ExpiresAt.IsZero() {
		return fmt.Errorf("validate signup OTP task: %w", asynq.SkipRetry)
	}

	if !time.Now().Before(payload.ExpiresAt) {
		return fmt.Errorf("signup OTP expired: %w", asynq.SkipRetry)
	}

	if err := w.ensureCurrentSignUpOTP(ctx, payload.ChallengeID, payload.OTPHash); err != nil {
		return err
	}

	if err := w.sender.SendSignUpOTP(ctx, payload.Email, payload.OTP); err != nil {
		return fmt.Errorf("send signup OTP: %w", err)
	}

	return nil
}

func (w *EmailWorker) ensureCurrentSignUpOTP(ctx context.Context, challengeID, otpHash string) error {
	challenge, err := auth.LoadSignupChallenge(ctx, w.rdb, challengeID)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return fmt.Errorf("signup challenge expired: %w", asynq.SkipRetry)
		}

		return fmt.Errorf("get signup challenge: %w", err)
	}

	if challenge.OTPHash != otpHash {
		return fmt.Errorf("signup OTP superseded: %w", asynq.SkipRetry)
	}

	return nil
}
