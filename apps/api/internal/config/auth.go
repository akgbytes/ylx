package config

import (
	"errors"
	"os"
	"strings"
	"time"
)

type AuthConfig struct {
	AccessTokenName            string
	RefreshTokenName           string
	AccessTokenExpiry          time.Duration
	RefreshTokenExpiry         time.Duration
	OTPLength                  int
	OTPMaxSends                int
	OTPMaxVerificationAttempts int
	OTPExpiry                  time.Duration
	OTPResendCooldown          time.Duration
	OTPSendLimitWindow         time.Duration
}

func loadAuthConfig() (AuthConfig, error) {
	otpLength, err := parseInt("OTP_LENGTH")
	if err != nil {
		return AuthConfig{}, err
	}

	otpMaxSends, err := parseInt("OTP_MAX_SENDS")
	if err != nil {
		return AuthConfig{}, err
	}

	otpMaxVerificationAttempts, err := parseInt("OTP_MAX_VERIFICATION_ATTEMPTS")
	if err != nil {
		return AuthConfig{}, err
	}

	accessTokenExpiry, err := parseDuration("ACCESS_TOKEN_EXPIRY")
	if err != nil {
		return AuthConfig{}, err
	}

	refreshTokenExpiry, err := parseDuration("REFRESH_TOKEN_EXPIRY")
	if err != nil {
		return AuthConfig{}, err
	}

	otpExpiry, err := parseDuration("OTP_EXPIRY")
	if err != nil {
		return AuthConfig{}, err
	}

	otpResendCooldown, err := parseDuration("OTP_RESEND_COOLDOWN")
	if err != nil {
		return AuthConfig{}, err
	}

	otpSendLimitWindow, err := parseDuration("OTP_SEND_LIMIT_WINDOW")
	if err != nil {
		return AuthConfig{}, err
	}

	return AuthConfig{
		AccessTokenName:            os.Getenv("ACCESS_TOKEN_NAME"),
		RefreshTokenName:           os.Getenv("REFRESH_TOKEN_NAME"),
		AccessTokenExpiry:          accessTokenExpiry,
		RefreshTokenExpiry:         refreshTokenExpiry,
		OTPLength:                  otpLength,
		OTPMaxSends:                otpMaxSends,
		OTPMaxVerificationAttempts: otpMaxVerificationAttempts,
		OTPExpiry:                  otpExpiry,
		OTPResendCooldown:          otpResendCooldown,
		OTPSendLimitWindow:         otpSendLimitWindow,
	}, nil
}

func (c *AuthConfig) validate() error {
	if c.AccessTokenName = strings.TrimSpace(c.AccessTokenName); c.AccessTokenName == "" {
		return errors.New("invalid configuration: ACCESS_TOKEN_NAME is required")
	}

	if c.RefreshTokenName = strings.TrimSpace(c.RefreshTokenName); c.RefreshTokenName == "" {
		return errors.New("invalid configuration: REFRESH_TOKEN_NAME is required")
	}

	if c.AccessTokenName == c.RefreshTokenName {
		return errors.New("invalid configuration: ACCESS_TOKEN_NAME and REFRESH_TOKEN_NAME must differ")
	}

	if c.AccessTokenExpiry <= 0 {
		return errors.New("invalid configuration: ACCESS_TOKEN_EXPIRY must be greater than 0")
	}

	if c.RefreshTokenExpiry <= 0 {
		return errors.New("invalid configuration: REFRESH_TOKEN_EXPIRY must be greater than 0")
	}

	if c.OTPLength <= 0 {
		return errors.New("invalid configuration: OTP_LENGTH must be greater than 0")
	}

	if c.OTPMaxSends <= 0 {
		return errors.New("invalid configuration: OTP_MAX_SENDS must be greater than 0")
	}

	if c.OTPMaxVerificationAttempts <= 0 {
		return errors.New("invalid configuration: OTP_MAX_VERIFICATION_ATTEMPTS must be greater than 0")
	}

	if c.OTPExpiry <= 0 {
		return errors.New("invalid configuration: OTP_EXPIRY must be greater than 0")
	}

	if c.OTPResendCooldown <= 0 {
		return errors.New("invalid configuration: OTP_RESEND_COOLDOWN must be greater than 0")
	}

	if c.OTPSendLimitWindow <= 0 {
		return errors.New("invalid configuration: OTP_SEND_LIMIT_WINDOW must be greater than 0")
	}

	return nil
}
