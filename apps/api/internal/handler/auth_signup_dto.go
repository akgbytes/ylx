package handler

import (
	"errors"
	"net/mail"
	"strings"
	"time"
)

type signupPayload struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type signUpResponse struct {
	ChallengeID string    `json:"challenge_id"`
	RetryAt     time.Time `json:"retry_at"`
}

type verifiedUserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type verifySignUpPayload struct {
	Email       string `json:"email"`
	ChallengeID string `json:"challenge_id"`
	OTP         string `json:"otp"`
}

type resendSignUpPayload struct {
	Email       string `json:"email"`
	ChallengeID string `json:"challenge_id"`
}

func (p *signupPayload) normalize() {
	p.Name = strings.TrimSpace(p.Name)
	p.Email = strings.TrimSpace(strings.ToLower(p.Email))
}

func (p *signupPayload) validate() (string, error) {
	if p.Name == "" {
		return "name", errors.New("name is required")
	}

	address, err := mail.ParseAddress(p.Email)
	// Need to check address bcuz ParseAddress() also passes 'Aman <aman@example.com>' string
	if err != nil || address.Address != p.Email {
		return "email", errors.New("invalid email address")
	}

	if len(p.Password) < 6 {
		return "password", errors.New("password must contain at least 6 characters")
	}

	return "", nil
}

func (p *verifySignUpPayload) normalize() {
	p.Email = strings.TrimSpace(strings.ToLower(p.Email))
	p.ChallengeID = strings.TrimSpace(p.ChallengeID)
	p.OTP = strings.TrimSpace(p.OTP)
}

func (p *verifySignUpPayload) validate() (string, error) {
	address, err := mail.ParseAddress(p.Email)
	if err != nil || address.Address != p.Email {
		return "email", errors.New("invalid email address")
	}

	if p.ChallengeID == "" {
		return "challenge_id", errors.New("challenge ID is required")
	}

	if p.OTP == "" {
		return "otp", errors.New("verification code is required")
	}

	return "", nil
}

func (p *resendSignUpPayload) normalize() {
	p.Email = strings.TrimSpace(strings.ToLower(p.Email))
	p.ChallengeID = strings.TrimSpace(p.ChallengeID)
}

func (p *resendSignUpPayload) validate() (string, error) {
	address, err := mail.ParseAddress(p.Email)
	if err != nil || address.Address != p.Email {
		return "email", errors.New("invalid email address")
	}

	if p.ChallengeID == "" {
		return "challenge_id", errors.New("challenge ID is required")
	}

	return "", nil
}
