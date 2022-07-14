package schema

import (
	"time"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/gateway/http/rest"
)

var (
	ErrMissingNameParameter = rest.ValidationError{
		Location: "body.name",
		Err:      rest.ErrMissingParameter,
	}

	ErrMissingCPFParameter = rest.ValidationError{
		Location: "body.cpf",
		Err:      rest.ErrMissingParameter,
	}

	ErrMissingSecretParameter = rest.ValidationError{
		Location: "body.secret",
		Err:      rest.ErrMissingParameter,
	}
)

type Account struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CPF       string `json:"cpf"`
	Balance   int    `json:"balance"`
	CreatedAt string `json:"created_at"`
}

type ListAccountsResponse struct {
	Accounts []Account `json:"accounts"`
}

func MapToListAccountsResponse(accs []entity.Account) ListAccountsResponse {
	response := ListAccountsResponse{
		Accounts: make([]Account, 0, len(accs)),
	}

	for _, acc := range accs {
		response.Accounts = append(response.Accounts, Account{
			ID:        acc.ID.String(),
			Name:      acc.Name,
			CPF:       acc.CPF.String(),
			Balance:   acc.Balance,
			CreatedAt: acc.CreatedAt.UTC().Format(time.RFC3339),
		})
	}

	return response
}

type BalanceResponse struct {
	Balance int `json:"balance"`
}

func MapToBalanceResponse(balance int) BalanceResponse {
	return BalanceResponse{Balance: balance}
}

type CreateAccountInput struct {
	Name   string `json:"name"`
	CPF    string `json:"cpf"`
	Secret string `json:"secret"`
}

func (r CreateAccountInput) IsValid() error {
	var errs rest.ValidationErrors

	if r.Name == "" {
		errs = append(errs, ErrMissingNameParameter)
	}

	if r.CPF == "" {
		errs = append(errs, ErrMissingCPFParameter)
	}

	if r.Secret == "" {
		errs = append(errs, ErrMissingSecretParameter)
	}

	if len(errs) != 0 {
		return errs
	}

	return nil
}

type CreateAccountResponse struct {
	Name      string `json:"name"`
	CPF       string `json:"cpf"`
	Balance   int    `json:"balance"`
	CreatedAt string `json:"created_at"`
}

func MapToCreateAccountResponse(acc entity.Account) CreateAccountResponse {
	return CreateAccountResponse{
		Name:      acc.Name,
		CPF:       acc.CPF.String(),
		Balance:   acc.Balance,
		CreatedAt: acc.CreatedAt.UTC().Format(time.RFC3339),
	}
}
