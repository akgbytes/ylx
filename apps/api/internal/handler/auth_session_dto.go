package handler

import (
	"errors"
	"net/mail"
	"strings"
)

type signInPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type signedInUserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (p *signInPayload) normalize() {
	p.Email = strings.TrimSpace(strings.ToLower(p.Email))
}

func (p *signInPayload) validate() (string, error) {
	address, err := mail.ParseAddress(p.Email)
	if err != nil || address.Address != p.Email {
		return "email", errors.New("invalid email address")
	}

	if p.Password == "" {
		return "password", errors.New("password is required")
	}

	return "", nil
}
