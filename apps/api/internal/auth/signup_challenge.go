package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/akgbytes/ylx/internal/config"
)

const (
	SignupReservationCooldownActive    = "cooldown_active"
	SignupReservationSendLimitReached  = "send_limit_reached"
	SignupReservationChallengeExpired  = "challenge_expired"
	SignupReservationChallengeMismatch = "challenge_mismatch"
)

type SignupChallenge struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	OTPHash      string `json:"otp_hash"`
}

type SignupReservation struct {
	Allowed bool
	Reason  string
	RetryAt time.Time
}

type signupReservationResponse struct {
	Allowed bool   `json:"allowed"`
	Reason  string `json:"reason"`
	RetryAt int64  `json:"retry_at"`
}

const signupReservationScript = `
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

const resendSignupChallengeScript = `
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
  local email = ARGV[5]
  local otpHash = ARGV[6]

  local challengeJSON = redis.call("GET", KEYS[3])
  if not challengeJSON then
    return response(false, "challenge_expired", 0)
  end

  local challenge = cjson.decode(challengeJSON)
  if challenge.email ~= email then
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

  return response(true, "ok", nowMs + cooldownMs)
`

func ReserveSignupChallenge(
	ctx context.Context,
	rdb *redis.Client,
	email, challengeID string,
	challenge SignupChallenge,
	authConfig config.AuthConfig,
) (SignupReservation, error) {
	challengeJSON, err := json.Marshal(challenge)
	if err != nil {
		return SignupReservation{}, fmt.Errorf("marshal signup challenge: %w", err)
	}

	rawResponse, err := redis.NewScript(signupReservationScript).Run(
		ctx,
		rdb,
		[]string{
			signupSendAttemptsKey(email),
			signupCooldownKey(email),
			signupChallengeKey(challengeID),
		},
		authConfig.OTPMaxSends,
		authConfig.OTPSendLimitWindow.Milliseconds(),
		authConfig.OTPResendCooldown.Milliseconds(),
		authConfig.OTPExpiry.Milliseconds(),
		string(challengeJSON),
	).Result()
	if err != nil {
		return SignupReservation{}, fmt.Errorf("reserve signup challenge: %w", err)
	}

	return decodeSignupReservation(rawResponse, "reserve signup challenge")
}

func ResendSignupChallenge(
	ctx context.Context,
	rdb *redis.Client,
	email, challengeID, otpHash string,
	authConfig config.AuthConfig,
) (SignupReservation, error) {
	rawResponse, err := redis.NewScript(resendSignupChallengeScript).Run(
		ctx,
		rdb,
		[]string{
			signupSendAttemptsKey(email),
			signupCooldownKey(email),
			signupChallengeKey(challengeID),
			signupVerifyAttemptsKey(challengeID),
		},
		authConfig.OTPMaxSends,
		authConfig.OTPSendLimitWindow.Milliseconds(),
		authConfig.OTPResendCooldown.Milliseconds(),
		authConfig.OTPExpiry.Milliseconds(),
		email,
		otpHash,
	).Result()
	if err != nil {
		return SignupReservation{}, fmt.Errorf("resend signup challenge: %w", err)
	}

	return decodeSignupReservation(rawResponse, "resend signup challenge")
}

func LoadSignupChallenge(ctx context.Context, rdb *redis.Client, challengeID string) (SignupChallenge, error) {
	data, err := rdb.Get(ctx, signupChallengeKey(challengeID)).Bytes()
	if err != nil {
		return SignupChallenge{}, fmt.Errorf("get signup challenge: %w", err)
	}

	var challenge SignupChallenge
	if err := json.Unmarshal(data, &challenge); err != nil {
		return SignupChallenge{}, fmt.Errorf("decode signup challenge: %w", err)
	}

	return challenge, nil
}

func decodeSignupReservation(rawResponse any, operation string) (SignupReservation, error) {
	jsonResponse, ok := rawResponse.(string)
	if !ok {
		return SignupReservation{}, fmt.Errorf("%s: unexpected Redis response type %T", operation, rawResponse)
	}

	var response signupReservationResponse
	if err := json.Unmarshal([]byte(jsonResponse), &response); err != nil {
		return SignupReservation{}, fmt.Errorf("decode %s response: %w", operation, err)
	}

	return SignupReservation{
		Allowed: response.Allowed,
		Reason:  response.Reason,
		RetryAt: time.UnixMilli(response.RetryAt),
	}, nil
}

func signupChallengeKey(challengeID string) string {
	return "auth:signup:challenge:" + challengeID
}

func signupVerifyAttemptsKey(challengeID string) string {
	return "auth:signup:verify-attempts:" + challengeID
}

func signupCooldownKey(email string) string {
	return "auth:signup:cooldown:" + email
}

func signupSendAttemptsKey(email string) string {
	return "auth:signup:send-attempts:" + email
}
