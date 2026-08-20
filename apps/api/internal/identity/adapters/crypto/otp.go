package crypto

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

func GenerateOTP() (string, error) {
	// Limit otp length to 6
	max := big.NewInt(1_000_000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	// Padding with 0 and limit length to 6
	return fmt.Sprintf("%06d", n), nil
}

func HashOTP(otp, secretKey string) string {
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(otp))
	return hex.EncodeToString(mac.Sum(nil))
}
