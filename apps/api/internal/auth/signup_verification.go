package auth

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type SignupVerificationAttemptResult int64

const (
	SignupVerificationExpired SignupVerificationAttemptResult = iota
	SignupVerificationRecorded
	SignupVerificationLimitReached
)

const recordFailedSignupVerificationScript = `
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

func RecordFailedSignupVerification(
	ctx context.Context,
	rdb *redis.Client,
	challengeID string,
	maxAttempts int,
) (SignupVerificationAttemptResult, error) {
	result, err := redis.NewScript(recordFailedSignupVerificationScript).Run(
		ctx,
		rdb,
		[]string{
			signupVerifyAttemptsKey(challengeID),
			signupChallengeKey(challengeID),
		},
		maxAttempts,
	).Int64()
	if err != nil {
		return SignupVerificationExpired, fmt.Errorf("record failed signup verification: %w", err)
	}

	return SignupVerificationAttemptResult(result), nil
}

func DeleteSignupVerificationState(ctx context.Context, rdb *redis.Client, challengeID string) error {
	if err := rdb.Del(
		ctx,
		signupChallengeKey(challengeID),
		signupVerifyAttemptsKey(challengeID),
	).Err(); err != nil {
		return fmt.Errorf("delete signup verification state: %w", err)
	}

	return nil
}
