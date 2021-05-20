package auth

import "errors"

var (
	errMissingCPFParameter    = errors.New("missing cpf parameter")
	errMissingSecretParameter = errors.New("missing secret parameter")
	errInvalidJSON            = errors.New("invalid json")
)

type responseBody struct {
	Token string `json:"token"`
}

type requestBody struct {
	CPF    string `json:"cpf"`
	Secret string `json:"secret"`
}

func (r requestBody) isValid() error {
	if r.CPF == "" {
		return errMissingCPFParameter
	}
	if r.Secret == "" {
		return errMissingSecretParameter
	}

	return nil
}
