package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserId    int
	Role      string
	SessionId string
	JTI       string
	jwt.RegisteredClaims
}

func GeneratateAccessToken(userId int, role string, sessionId string, tokenKey string) (string, error) {

	claims := Claims{
		UserId:    userId,
		Role:      role,
		SessionId: sessionId,
		JTI:       uuid.NewString(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "bank",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(tokenKey))

	if err != nil {
		return "", fmt.Errorf("sign token : %w", err)
	}

	return tokenString, nil

}

func GeneratateRefreshToken(userId int, role string, sessionId string, tokenKey string) (string, *time.Time, *time.Time, error) {
	expiresAt := jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour))
	issuedAt := jwt.NewNumericDate(time.Now())

	claims := Claims{
		UserId:    userId,
		Role:      role,
		SessionId: sessionId,
		JTI:       uuid.NewString(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expiresAt,
			IssuedAt:  issuedAt,
			Issuer:    "bank",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(tokenKey))

	if err != nil {
		return "", nil, nil, fmt.Errorf("sign token : %w", err)
	}

	return tokenString, &expiresAt.Time, &issuedAt.Time, nil

}
func ParseAndValidateToken(tokenString string, tokenKey string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error in parse token:")
		}

		return []byte(tokenKey), nil

	})

	if err != nil {
		return nil, fmt.Errorf("jwt parser: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func HashRefreshToken(tokenString string, tokenKey string) string {
	mac := hmac.New(sha256.New, []byte(tokenKey))
	mac.Write([]byte(tokenString))
	return hex.EncodeToString(mac.Sum(nil))
}
