package crypto

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashEmail(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}
