package vos

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"

	"github.com/Nhanderu/brdoc"
)

var regexOnlyDigitsCPF = regexp.MustCompile(`^(\d{3})(\d{3})(\d{3})(\d{2})$`)

// ErrInvalidCPF occurs when the cpf is not valid or is not formatted correctly.
var ErrInvalidCPF = errors.New("invalid cpf")

type CPF struct {
	value string
}

// IsValid checks if the CPF is valid.
// For now, I decided to use a third-party library and, later on,
// create my own cpf validation implementation.
func (c *CPF) IsValid() bool {
	if ok := brdoc.IsCPF(c.value); !ok {
		return false
	}

	// If the CPF is valid, it means that it's formatted as XXX.XXX.XXX-XX
	// or XXXXXXXXXXX. The condition below guarantees that the cpf will be saved
	// without any pontcuation.
	if onlyDigits := regexOnlyDigitsCPF.MatchString(c.value); !onlyDigits {
		reg := regexp.MustCompile("[^0-9]+")
		c.value = reg.ReplaceAllString(c.value, "")
	}

	return true
}

func (c CPF) String() string {
	return c.value
}

func NewCPF(cpf string) (CPF, error) {
	c := CPF{value: cpf}
	if !c.IsValid() {
		return CPF{}, ErrInvalidCPF
	}

	return c, nil
}

func (c CPF) Value() (driver.Value, error) {
	return c.String(), nil
}

func (c *CPF) Scan(value interface{}) error {
	if value == nil {
		*c = CPF{}

		return nil
	}

	bv, err := driver.String.ConvertValue(value)
	if err == nil {
		if v, ok := bv.(string); ok {
			*c = CPF{v}

			return nil
		}
	}

	return fmt.Errorf("could not scan cpf: %w", err)
}
