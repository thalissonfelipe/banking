package schemes

import "errors"

// Request input errors.
var (
	ErrMissingCPFParameter    = errors.New("missing cpf parameter")
	ErrMissingSecretParameter = errors.New("missing secret parameter")
)

type LoginResponse struct {
	Token string `json:"token"`
}

type LoginInput struct {
	CPF    string `json:"cpf"`
	Secret string `json:"secret"`
}

func (r LoginInput) IsValid() error {
	if r.CPF == "" {
		return ErrMissingCPFParameter
	}

	if r.Secret == "" {
		return ErrMissingSecretParameter
	}

	return nil
}
