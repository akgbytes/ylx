package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"golang.org/x/crypto/argon2"
)

type ArgonParams struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

func DefaultArgonParams() *ArgonParams {
	return &ArgonParams{
		Memory:      64 * 1024, // 64 MB
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  12,
		KeyLength:   32,
	}
}

func GenerateOTP() (string, error) {
	// Setting upper bound to limit otp length to 6
	otpUpperBound := big.NewInt(1_000_000)
	n, err := rand.Int(rand.Reader, otpUpperBound)
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

func GenerateRandomBytes(n uint32) ([]byte, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}

	return bytes, nil
}

func GenerateRandomString(n uint32) (string, error) {
	bytes, err := GenerateRandomBytes(n)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func HashPassword(password string, p *ArgonParams) (string, error) {
	salt, err := GenerateRandomBytes(p.SaltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		p.Memory,
		p.Iterations,
		p.Parallelism,
		b64Salt,
		b64Hash,
	)

	return encoded, nil
}

func VerifyPassword(password, passwordHash string) (bool, error) {
	p, salt, hash, err := decodeHash(passwordHash)
	if err != nil {
		return false, err
	}

	otherHash := argon2.IDKey([]byte(password), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}

	return false, nil
}

func decodeHash(encodedHash string) (*ArgonParams, []byte, []byte, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return nil, nil, nil, errors.New("invalid hash format")
	}

	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, errors.New("incompatible argon2 version")
	}

	p := &ArgonParams{}
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &p.Memory, &p.Iterations, &p.Parallelism); err != nil {
		return nil, nil, nil, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, err
	}

	// Checking if salt is greater than maximum value of uint32
	maxUint32 := ^uint32(0)
	if uint64(len(salt)) > uint64(maxUint32) {
		return nil, nil, nil, errors.New("salt is too long")
	}
	p.SaltLength = uint32(len(salt)) // #nosec G115 -- length is checked against maxUint32 above.

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, nil, err
	}

	if uint64(len(hash)) > uint64(maxUint32) {
		return nil, nil, nil, errors.New("hash is too long")
	}
	p.KeyLength = uint32(len(hash)) // #nosec G115 -- length is checked against maxUint32 above.

	return p, salt, hash, nil
}
