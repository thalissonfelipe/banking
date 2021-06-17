package schemes

import "errors"

// Request input errors.
var (
	ErrMissingNameParameter   = errors.New("missing name parameter")
	ErrMissingCPFParameter    = errors.New("missing cpf parameter")
	ErrMissingSecretParameter = errors.New("missing secret parameter")
)

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
		return ErrMissingNameParameter
	}

	if r.CPF == "" {
		return ErrMissingCPFParameter
	}

	if r.Secret == "" {
		return ErrMissingSecretParameter
	}

	return nil
}
