package generateotp

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

func GenerateOtp() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(900000))

	if err != nil {
		return "", fmt.Errorf("Error in generationg otp: %w", err)
	}

	return fmt.Sprintf("%d", n.Int64()+100000), nil
}

func HashOtp(otp string) (string, error) {
	if otp == "" {
		return "", fmt.Errorf("Otp is required")
	}

	hash := sha256.Sum256([]byte(otp))

	return hex.EncodeToString(hash[:]), nil
}
