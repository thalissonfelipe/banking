package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

var jwtKey = []byte("secret-key")

type claims struct {
	AccountOriginID string `json:"account_origin_id"`
	jwt.StandardClaims
}

func NewToken(accountOriginID string) (string, error) {
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

func IsValidToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		log.WithError(err).Error("unexpected error ocurred during parse jwt")
		return ErrUnauthorized
	}
	if !token.Valid {
		return ErrUnauthorized
	}

	return nil
}

func GetIDFromToken(tokenString string) string {
	claims := &claims{}
	jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	return claims.AccountOriginID
}
