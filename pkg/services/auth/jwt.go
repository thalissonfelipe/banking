package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

var jwtKey = []byte("secret-key")

type claims struct {
	AccountOriginID string `json:"account_origin_id"`
	jwt.StandardClaims
}

func newToken(accountOriginID string) (string, error) {
	expirationTime := time.Now().Add(10 * time.Minute)
	claims := &claims{
		AccountOriginID: accountOriginID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", entities.ErrInternalError
	}

	return tokenString, nil
}

func GetIDFromToken(tokenString string) (string, error) {
	claims := &claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return "", ErrUnauthorized
	}

	return claims.AccountOriginID, nil
}
