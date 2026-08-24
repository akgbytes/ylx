package otpstore

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const resendChallengeScript = `
local function response(allowed, reason, retryAtMs, recipient, email)
  return cjson.encode({
    allowed = allowed,
    reason = reason,
    retry_at = retryAtMs,
    recipient = recipient or "",
    email = email or ""
  })
end

local now = redis.call("TIME")
local nowMs = tonumber(now[1]) * 1000 +
  math.floor(tonumber(now[2]) / 1000)

local maxSends = tonumber(ARGV[1])
local sendLimitWindowMs = tonumber(ARGV[2])
local cooldownMs = tonumber(ARGV[3])
local challengeTTLms = tonumber(ARGV[4])
local emailHash = ARGV[5]
local otpHash = ARGV[6]

local challengeJSON = redis.call("GET", KEYS[3])
if not challengeJSON then
  return response(false, "challenge_expired", 0)
end

local challenge = cjson.decode(challengeJSON)
if challenge.email_hash ~= emailHash then
  return response(false, "challenge_mismatch", 0)
end

local attempts = tonumber(redis.call("GET", KEYS[1]) or "0")
local attemptsTTL = redis.call("PTTL", KEYS[1])

if attempts > 0 and attemptsTTL <= 0 then
  return response(false, "invalid_attempts_state", 0)
end

if attempts >= maxSends then
  return response(false, "send_limit_reached", nowMs + attemptsTTL)
end

local cooldownTTL = redis.call("PTTL", KEYS[2])
if cooldownTTL == -1 then
  return response(false, "invalid_cooldown_state", 0)
end

if cooldownTTL > 0 then
  return response(false, "cooldown_active", nowMs + cooldownTTL)
end

challenge.otp_hash = otpHash
redis.call("SET", KEYS[3], cjson.encode(challenge), "PX", challengeTTLms)
redis.call("DEL", KEYS[4])
redis.call("SET", KEYS[2], "1", "PX", cooldownMs)

local newAttempts = redis.call("INCR", KEYS[1])
if newAttempts == 1 then
  redis.call("PEXPIRE", KEYS[1], sendLimitWindowMs)
end

return response(true, "ok", nowMs + cooldownMs, challenge.name, challenge.email)
`

func (s *Store) Resend(ctx context.Context, emailHash, otpHash string) (ResendResult, error) {
	rawResponse, err := redis.NewScript(resendChallengeScript).Run(
		ctx,
		s.rdb,
		[]string{
			signupSendAttemptsKey(emailHash),
			signupCooldownKey(emailHash),
			signupChallengeKey(emailHash),
			signupVerificationAttemptsKey(emailHash),
		},
		s.cfg.OTPMaxSends,
		s.cfg.OTPSendLimitWindow.Milliseconds(),
		s.cfg.OTPResendCooldown.Milliseconds(),
		s.cfg.OTPExpiry.Milliseconds(),
		emailHash,
		otpHash,
	).Result()
	if err != nil {
		return ResendResult{}, fmt.Errorf("resend signup challenge: %w", err)
	}

	jsonResponse, ok := rawResponse.(string)
	if !ok {
		return ResendResult{}, fmt.Errorf("resend signup challenge: unexpected redis response type %T", rawResponse)
	}

	var response struct {
		Allowed   bool   `json:"allowed"`
		Reason    string `json:"reason"`
		RetryAt   int64  `json:"retry_at"`
		ExpiresAt int64  `json:"expires_at"`
		Recipient string `json:"recipient"`
		Email     string `json:"email"`
	}
	if err := json.Unmarshal([]byte(jsonResponse), &response); err != nil {
		return ResendResult{}, fmt.Errorf("decode resend signup challenge response: %w", err)
	}

	reservation := Reservation{
		Allowed: response.Allowed,
		Reason:  response.Reason,
		RetryAt: time.UnixMilli(response.RetryAt),
	}
	if response.ExpiresAt > 0 {
		reservation.ExpiresAt = time.UnixMilli(response.ExpiresAt)
	}

	return ResendResult{
		Reservation: reservation,
		Recipient:   response.Recipient,
		Email:       response.Email,
	}, nil
}
