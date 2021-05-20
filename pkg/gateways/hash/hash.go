package hash

import "golang.org/x/crypto/bcrypt"

type Hash struct{}

func (h Hash) Hash(secret string) ([]byte, error) {
	hashedSecret, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	return hashedSecret, err
}

func (h Hash) CompareHashAndSecret(hashedSecret, secret []byte) error {
	return bcrypt.CompareHashAndPassword(hashedSecret, secret)
}
