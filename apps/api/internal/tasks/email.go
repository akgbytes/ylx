package tasks

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)

const TypeSendSignUpOTP = "email:send_signup:otp"

type SendSignUpOTPPayload struct {
	ChallengeID string    `json:"challenge_id"`
	Email       string    `json:"email"`
	OTP         string    `json:"otp"`
	OTPHash     string    `json:"otp_hash"`
	ExpiresAt   time.Time `json:"expires_at"`
}

func NewSendSignUpOTPTask(p SendSignUpOTPPayload) (*asynq.Task, error) {
	data, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("marshal signup OTP task: %w", err)
	}

	return asynq.NewTask(TypeSendSignUpOTP, data), nil
}
