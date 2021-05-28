package vos

import (
	"database/sql/driver"
	"errors"

	"github.com/Nhanderu/brdoc"
)

var ErrInvalidCPF = errors.New("invalid cpf")

type CPF struct {
	value string
}

// IsValid checks if the CPF is valid.
// For now, I decided to use a third-party library and, later on,
// create my own cpf validation implementation.
func (c CPF) IsValid() bool {
	return brdoc.IsCPF(c.value)
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
