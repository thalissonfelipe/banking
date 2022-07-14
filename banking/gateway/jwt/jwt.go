package jwt

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/thalissonfelipe/banking/banking/domain/usecases"
)

var jwtKey = []byte("secret-key")

type JWT struct{}

func New() *JWT {
	return &JWT{}
}

func (JWT) NewToken(accountID string) (string, error) {
	return NewToken(accountID)
}

type Claims struct {
	AccountID string `json:"account_id"`
	jwt.RegisteredClaims
}

func NewToken(accountID string) (string, error) {
	const expiresAt = 24 * time.Hour // 1 day

	claims := &Claims{
		AccountID: accountID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresAt)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		return "", fmt.Errorf("getting signed token: %w", err)
	}

	return signedToken, nil
}

func IsTokenValid(tokenStr string) error {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method: %v: %w", t.Header["alg"], usecases.ErrUnauthorized)
		}

		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return usecases.ErrUnauthorized
	}

	return nil
}

func GetAccountIDFromToken(tokenStr string) string {
	claims := Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method: %v: %w", t.Header["alg"], usecases.ErrUnauthorized)
		}

		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		log.Println(err)
		return ""
	}

	return claims.AccountID
}
