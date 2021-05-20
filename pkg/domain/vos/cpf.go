package vos

import "github.com/Nhanderu/brdoc"

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

func NewCPF(cpf string) CPF {
	return CPF{value: cpf}
}
