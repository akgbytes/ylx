package rest

import (
	"errors"
	"net/mail"
	"strings"
	"time"

	"github.com/akgbytes/ylx/internal/identity/internal/domain"
)

type signupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (p *signupRequest) normalize() {
	p.Name = strings.TrimSpace(p.Name)
	p.Email = strings.TrimSpace(strings.ToLower(p.Email))
}

func (p *signupRequest) validate() (string, error) {
	if p.Name == "" {
		return "name", errors.New("name is required")
	}

	if len(p.Password) < 6 {
		return "password", errors.New("password must contain at least 6 characters")
	}

	// ParseAddress also accepts display names (e.g. "Aman <akgbytes@gmail.com>")
	// So reject addresses with display names
	addr, err := mail.ParseAddress(p.Email)
	if err != nil || addr.Address != p.Email {
		return "email", errors.New("invalid email address")
	}

	return "", nil
}

type signupResponse struct {
	RetryAt time.Time `json:"retry_at"`
}

type resendSignupRequest struct {
	Email string `json:"email"`
}

func (p *resendSignupRequest) normalize() {
	p.Email = strings.TrimSpace(strings.ToLower(p.Email))
}

func (p *resendSignupRequest) validate() (string, error) {
	addr, err := mail.ParseAddress(p.Email)
	if err != nil || addr.Address != p.Email {
		return "email", errors.New("invalid email address")
	}
	return "", nil
}

type verifySignupRequest struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

func (p *verifySignupRequest) normalize() {
	p.Email = strings.TrimSpace(strings.ToLower(p.Email))
	p.OTP = strings.TrimSpace(p.OTP)
}

func (p *verifySignupRequest) validate() (string, error) {
	addr, err := mail.ParseAddress(p.Email)
	if err != nil || addr.Address != p.Email {
		return "email", errors.New("invalid email address")
	}

	if p.OTP == "" {
		return "otp", errors.New("verification code is required")
	}

	return "", nil
}

type userResponse struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

func newUserResponse(user domain.User) userResponse {
	return userResponse{
		ID:        user.ID.String(),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: &user.CreatedAt,
	}
}

type signinRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (p *signinRequest) normalize() {
	p.Email = strings.TrimSpace(strings.ToLower(p.Email))
}

func (p *signinRequest) validate() (string, error) {
	addr, err := mail.ParseAddress(p.Email)
	if err != nil || addr.Address != p.Email {
		return "email", errors.New("invalid email address")
	}

	if p.Password == "" {
		return "password", errors.New("password is required")
	}

	return "", nil
}
