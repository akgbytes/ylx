package rest

import (
	"errors"
	"net/mail"
	"strings"
	"time"
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
