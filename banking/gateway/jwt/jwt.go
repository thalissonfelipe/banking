package jwt

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/thalissonfelipe/banking/banking/domain/usecases"
)

var (
	jwtKey         = []byte("secret-key")
	expirationTime = 24 * time.Hour
)

type Claims struct {
	AccountOriginID string `json:"account_origin_id"`
	jwt.StandardClaims
}

func NewToken(accountOriginID string) (string, error) {
	expirationTime := time.Now().Add(expirationTime)
	claims := &Claims{
		AccountOriginID: accountOriginID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", fmt.Errorf("getting signed token: %w", err)
	}

	return tokenString, nil
}

func IsValidToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return usecases.ErrUnauthorized
	}

	return nil
}

func GetIDFromToken(tokenString string) string {
	claims := &Claims{}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return ""
	}

	return claims.AccountOriginID
}
