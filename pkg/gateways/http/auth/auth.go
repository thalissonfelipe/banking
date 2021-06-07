package auth

import "errors"

var (
	errMissingCPFParameter    = errors.New("missing cpf parameter")
	errMissingSecretParameter = errors.New("missing secret parameter")
)

type LoginResponse struct {
	Token string `json:"token"`
}

type LoginInput struct {
	CPF    string `json:"cpf"`
	Secret string `json:"secret"`
}

func (r LoginInput) isValid() error {
	if r.CPF == "" {
		return errMissingCPFParameter
	}

	if r.Secret == "" {
		return errMissingSecretParameter
	}

	return nil
}
