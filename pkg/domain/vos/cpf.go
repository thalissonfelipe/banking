package vos

import (
	"database/sql/driver"
	"errors"
	"log"
	"regexp"

	"github.com/Nhanderu/brdoc"
)

var regexOnlyDigitsCPF = regexp.MustCompile(`^(\d{3})(\d{3})(\d{3})(\d{2})$`)

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
		reg, err := regexp.Compile("[^0-9]+")
		if err != nil {
			log.Fatal(err)
		}

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
		*c = CPF(CPF{})
		return nil
	}
	if bv, err := driver.String.ConvertValue(value); err == nil {
		if v, ok := bv.(string); ok {
			*c = CPF(CPF{v})
			return nil
		}
	}

	return errors.New("failed to scan CPF")
}
