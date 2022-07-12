package vos

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"

	"github.com/thalissonfelipe/banking/banking/domain/encrypter"
)

// ErrInvalidSecret occurs when the secret does meet the basic requirements. Also, the secret
// must have a minimum size of 8 and a maximum size of 20 characters. Must have at least one
// uppercase characterer, at least one lowercase character and at least one number.
var ErrInvalidSecret = errors.New("invalid secret")

const (
	secretMaxSize = 20
	secretMinSize = 8
)

var (
	regexUpper  = regexp.MustCompile("[A-Z]")
	regexLower  = regexp.MustCompile("[a-z]")
	regexNumber = regexp.MustCompile("[0-9]")
)

type Secret struct {
	value string
}

func (s Secret) IsValid() bool {
	if s.Size() < secretMinSize {
		return false
	}

	if s.Size() > secretMaxSize {
		return false
	}

	if hasUpper := regexUpper.MatchString(s.value); !hasUpper {
		return false
	}

	if hasLower := regexLower.MatchString(s.value); !hasLower {
		return false
	}

	if hasNumber := regexNumber.MatchString(s.value); !hasNumber {
		return false
	}

	return true
}

func (s Secret) String() string {
	return s.value
}

func (s Secret) Size() int {
	return len(s.value)
}

func (s Secret) Value() (driver.Value, error) {
	return s.String(), nil
}

func (s *Secret) Scan(value interface{}) error {
	if value == nil {
		*s = Secret{}

		return nil
	}

	bv, err := driver.String.ConvertValue(value)
	if err == nil {
		if v, ok := bv.(string); ok {
			*s = Secret{v}

			return nil
		}
	}

	return fmt.Errorf("could not scan secret: %w", err)
}

func (s *Secret) Hash(encrypter encrypter.Encrypter) error {
	hashedSecret, err := encrypter.Hash(s.value)
	if err != nil {
		return fmt.Errorf("could not encrypt secret: %w", err)
	}

	s.value = string(hashedSecret)

	return nil
}

func NewSecret(secret string) (Secret, error) {
	s := Secret{value: secret}
	if !s.IsValid() {
		return Secret{}, ErrInvalidSecret
	}

	return s, nil
}
