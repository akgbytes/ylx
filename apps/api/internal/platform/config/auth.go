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
	JWTSecretKey               string
	OTPSecretKey               string
	OTPMaxSends                int
	OTPMaxVerificationAttempts int
	OTPExpiry                  time.Duration
	OTPResendCooldown          time.Duration
	OTPSendLimitWindow         time.Duration
}

func loadAuthConfig() (AuthConfig, error) {
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
		JWTSecretKey:               os.Getenv("JWT_SECRET_KEY"),
		OTPSecretKey:               os.Getenv("OTP_SECRET_KEY"),
		OTPMaxSends:                otpMaxSends,
		OTPMaxVerificationAttempts: otpMaxVerificationAttempts,
		OTPExpiry:                  otpExpiry,
		OTPResendCooldown:          otpResendCooldown,
		OTPSendLimitWindow:         otpSendLimitWindow,
	}, nil
}

func (c *AuthConfig) normalize() {
	c.AccessTokenName = strings.TrimSpace(c.AccessTokenName)
	c.RefreshTokenName = strings.TrimSpace(c.RefreshTokenName)
	c.JWTSecretKey = strings.TrimSpace(c.JWTSecretKey)
	c.OTPSecretKey = strings.TrimSpace(c.OTPSecretKey)
}

func (c *AuthConfig) validate() error {
	if c.AccessTokenName == "" {
		return errors.New("invalid configuration: ACCESS_TOKEN_NAME is required")
	}

	if c.RefreshTokenName == "" {
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

	if len(c.JWTSecretKey) < 32 {
		return errors.New("invalid configuration: JWT_SECRET_KEY must be at least 32 characters")
	}

	if c.OTPSecretKey == "" {
		return errors.New("invalid configuration: OTP_SECRET_KEY is required")
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
