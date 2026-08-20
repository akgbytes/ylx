package otpstore

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type AttemptResult int64

const (
	AttemptExpired AttemptResult = iota
	AttemptRecorded
	AttemptLimitReached
)

const recordFailedAttemptScript = `
local challengeTTL = redis.call("PTTL", KEYS[2])

if challengeTTL <= 0 then
  return 0
end

local attempts = redis.call("INCR", KEYS[1])

if attempts == 1 then
  redis.call("PEXPIRE", KEYS[1], challengeTTL)
end

if attempts >= tonumber(ARGV[1]) then
  redis.call("DEL", KEYS[1])
  redis.call("DEL", KEYS[2])
  return 2
end

return 1
`

func (s *Store) RecordFailedAttempt(ctx context.Context, emailHash string) (AttemptResult, error) {
	result, err := redis.NewScript(recordFailedAttemptScript).Run(
		ctx,
		s.rdb,
		[]string{
			signupVerificationAttemptsKey(emailHash),
			signupChallengeKey(emailHash),
		},
		s.cfg.OTPMaxVerificationAttempts,
	).Int64()
	if err != nil {
		return AttemptExpired, fmt.Errorf("record failed signup verification: %w", err)
	}

	return AttemptResult(result), nil
}

func (s *Store) Clear(ctx context.Context, emailHash string) error {
	if err := s.rdb.Del(
		ctx,
		signupChallengeKey(emailHash),
		signupVerificationAttemptsKey(emailHash),
	).Err(); err != nil {
		return fmt.Errorf("delete signup verification state: %w", err)
	}

	return nil
}
