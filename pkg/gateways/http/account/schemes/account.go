package schemes

import "github.com/thalissonfelipe/banking/pkg/gateways/http/responses"

type AccountListResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CPF       string `json:"cpf"`
	Balance   int    `json:"balance"`
	CreatedAt string `json:"created_at"`
}

type BalanceResponse struct {
	Balance int `json:"balance"`
}

type CreateAccountInput struct {
	Name   string `json:"name"`
	CPF    string `json:"cpf"`
	Secret string `json:"secret"`
}

type CreateAccountResponse struct {
	Name    string `json:"name"`
	CPF     string `json:"cpf"`
	Balance int    `json:"balance"`
}

func (r CreateAccountInput) IsValid() error {
	if r.Name == "" {
		return responses.ErrMissingNameParameter
	}

	if r.CPF == "" {
		return responses.ErrMissingCPFParameter
	}

	if r.Secret == "" {
		return responses.ErrMissingSecretParameter
	}

	return nil
}
