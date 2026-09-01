package task

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"

	"github.com/akgbytes/ylx/internal/identity/internal/app"
)

type Dispatcher struct {
	client *asynq.Client
}

func NewDispatcher(client *asynq.Client) *Dispatcher {
	return &Dispatcher{client: client}
}

func (d *Dispatcher) DispatchSignupOTP(ctx context.Context, message app.SignupOTP) error {
	payload := signupOTPTaskPayload{
		Recipient: message.Recipient,
		Email:     message.Email,
		EmailHash: message.EmailHash,
		OTP:       message.OTP,
		OTPHash:   message.OTPHash,
		ExpiresAt: message.ExpiresAt,
	}

	task, err := newSignupOTPTask(payload)
	if err != nil {
		return err
	}

	if _, err := d.client.EnqueueContext(
		ctx,
		task,
		asynq.Queue(Queue),
		asynq.MaxRetry(maxRetry),
		asynq.Timeout(taskTimeout),
		asynq.Deadline(payload.ExpiresAt),
	); err != nil {
		return fmt.Errorf("enqueue signup otp task: %w", err)
	}

	return nil
}
