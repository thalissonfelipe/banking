package schemes

import "github.com/thalissonfelipe/banking/pkg/gateways/http/responses"

type LoginResponse struct {
	Token string `json:"token"`
}

type LoginInput struct {
	CPF    string `json:"cpf"`
	Secret string `json:"secret"`
}

func (r LoginInput) IsValid() error {
	if r.CPF == "" {
		return responses.ErrMissingCPFParameter
	}

	if r.Secret == "" {
		return responses.ErrMissingSecretParameter
	}

	return nil
}
