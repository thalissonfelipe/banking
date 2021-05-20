package vos

import (
	"database/sql/driver"
	"errors"
	"regexp"
)

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

func NewSecret(secret string) Secret {
	return Secret{value: secret}
}

func (c Secret) Value() (driver.Value, error) {
	return c.String(), nil
}

func (c *Secret) Scan(value interface{}) error {
	if value == nil {
		*c = Secret(Secret{})
		return nil
	}
	if bv, err := driver.String.ConvertValue(value); err == nil {
		if v, ok := bv.(string); ok {
			*c = Secret(Secret{v})
			return nil
		}
	}

	return errors.New("failed to scan Secret")
}
