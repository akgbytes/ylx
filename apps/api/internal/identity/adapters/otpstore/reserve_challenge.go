package otpstore

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Challenge struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	EmailHash    string `json:"email_hash"`
	PasswordHash string `json:"password_hash"`
	OTPHash      string `json:"otp_hash"`
}

type Reservation struct {
	Allowed bool
	Reason  string
	RetryAt time.Time
}

type reservationResponse struct {
	Allowed bool   `json:"allowed"`
	Reason  string `json:"reason"`
	RetryAt int64  `json:"retry_at"`
}

const (
	ReasonCooldownActive    = "cooldown_active"
	ReasonSendLimitReached  = "send_limit_reached"
	ReasonChallengeExpired  = "challenge_expired"
	ReasonChallengeMismatch = "challenge_mismatch"
)

const reserveScript = `
local function response(allowed, reason, retryAtMs)
  return cjson.encode({
    allowed = allowed,
    reason = reason,
    retry_at = retryAtMs
  })
end

local now = redis.call("TIME")
local nowMs = tonumber(now[1]) * 1000 +
  math.floor(tonumber(now[2]) / 1000)

local maxSends = tonumber(ARGV[1])
local sendLimitWindowMs = tonumber(ARGV[2])
local cooldownMs = tonumber(ARGV[3])
local challengeTTLms = tonumber(ARGV[4])
local challengeJSON = ARGV[5]

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

redis.call("SET", KEYS[3], challengeJSON, "PX", challengeTTLms)
redis.call("SET", KEYS[2], "1", "PX", cooldownMs)

local newAttempts = redis.call("INCR", KEYS[1])

if newAttempts == 1 then
  redis.call("PEXPIRE", KEYS[1], sendLimitWindowMs)
end

return response(true, "ok", nowMs + cooldownMs)
`

func (s *Store) Reserve(ctx context.Context, challenge Challenge) (Reservation, error) {
	challengeJSON, err := json.Marshal(challenge)
	if err != nil {
		return Reservation{}, fmt.Errorf("marshal signup challenge: %w", err)
	}

	rawResponse, err := redis.NewScript(reserveScript).Run(
		ctx,
		s.rdb,
		[]string{
			signupSendAttemptsKey(challenge.EmailHash),
			signupCooldownKey(challenge.EmailHash),
			signupChallengeKey(challenge.EmailHash),
		},
		s.cfg.OTPMaxSends,
		s.cfg.OTPSendLimitWindow.Milliseconds(),
		s.cfg.OTPResendCooldown.Milliseconds(),
		s.cfg.OTPExpiry.Milliseconds(),
		string(challengeJSON),
	).Result()
	if err != nil {
		return Reservation{}, fmt.Errorf("reserve signup challenge: %w", err)
	}

	return decodeReservation(rawResponse, "reserve signup challenge")
}

func (s *Store) Load(ctx context.Context, emailHash string) (Challenge, error) {
	data, err := s.rdb.Get(ctx, signupChallengeKey(emailHash)).Bytes()
	if err != nil {
		return Challenge{}, fmt.Errorf("get signup challenge: %w", err)
	}

	var challenge Challenge
	if err := json.Unmarshal(data, &challenge); err != nil {
		return Challenge{}, fmt.Errorf("decode signup challenge: %w", err)
	}

	return challenge, nil
}

func decodeReservation(rawResponse any, operation string) (Reservation, error) {
	jsonResponse, ok := rawResponse.(string)
	if !ok {
		return Reservation{}, fmt.Errorf("%s: unexpected redis response type %T", operation, rawResponse)
	}

	var response reservationResponse
	if err := json.Unmarshal([]byte(jsonResponse), &response); err != nil {
		return Reservation{}, fmt.Errorf("decode %s response: %w", operation, err)
	}

	return Reservation{
		Allowed: response.Allowed,
		Reason:  response.Reason,
		RetryAt: time.UnixMilli(response.RetryAt),
	}, nil
}
