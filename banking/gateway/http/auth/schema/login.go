package schema

import "github.com/thalissonfelipe/banking/banking/gateway/http/rest"

type LoginResponse struct {
	Token string `json:"token"`
}

func MapToLoginResponse(token string) LoginResponse {
	return LoginResponse{Token: token}
}

type LoginInput struct {
	CPF    string `json:"cpf"`
	Secret string `json:"secret"`
}

func (r LoginInput) IsValid() error {
	if r.CPF == "" {
		return rest.ErrMissingCPFParameter
	}

	if r.Secret == "" {
		return rest.ErrMissingSecretParameter
	}

	return nil
}