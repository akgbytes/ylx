package task

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)

const (
	TypeSignupOTP = "email:signup_otp"
	Queue         = "email"
	maxRetry      = 5
	taskTimeout   = 15 * time.Second
)

type signupOTPTaskPayload struct {
	Recipient string    `json:"recipient"`
	Email     string    `json:"email"`
	EmailHash string    `json:"email_hash"`
	OTP       string    `json:"otp"`
	OTPHash   string    `json:"otp_hash"`
	ExpiresAt time.Time `json:"expires_at"`
}

func newSignupOTPTask(payload signupOTPTaskPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal signup otp task: %w", err)
	}
	return asynq.NewTask(TypeSignupOTP, data), nil
}
