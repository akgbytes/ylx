package config

import (
	"errors"
	"net/mail"
	"os"
	"strings"
)

type EmailConfig struct {
	ResendAPIKey string
	From         string
}

func loadEmailConfig() EmailConfig {
	return EmailConfig{
		ResendAPIKey: os.Getenv("RESEND_API_KEY"),
		From:         os.Getenv("EMAIL_FROM"),
	}
}

func (c *EmailConfig) normalize() {
	c.ResendAPIKey = strings.TrimSpace(c.ResendAPIKey)
	c.From = strings.TrimSpace(c.From)
}

func (c *EmailConfig) validate() error {
	if c.ResendAPIKey == "" {
		return errors.New("invalid configuration: RESEND_API_KEY is required")
	}

	address, err := mail.ParseAddress(c.From)
	if err != nil || address.Address == "" {
		return errors.New("invalid configuration: EMAIL_FROM must be a valid email address")
	}

	return nil
}
