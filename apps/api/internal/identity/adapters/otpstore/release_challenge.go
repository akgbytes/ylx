package otpstore

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

const releaseScript = `
local attempts = tonumber(redis.call("GET", KEYS[1]) or "0")

if attempts > 0 then
  redis.call("DECR", KEYS[1])
end

redis.call("DEL", KEYS[2])
redis.call("DEL", KEYS[3])
redis.call("DEL", KEYS[4])

return 1
`

func (s *Store) Release(ctx context.Context, emailHash string) error {
	if err := redis.NewScript(releaseScript).Run(
		ctx,
		s.rdb,
		[]string{
			signupSendAttemptsKey(emailHash),
			signupCooldownKey(emailHash),
			signupChallengeKey(emailHash),
			signupVerificationAttemptsKey(emailHash),
		},
	).Err(); err != nil {
		return fmt.Errorf("release signup reservation: %w", err)
	}

	return nil
}
