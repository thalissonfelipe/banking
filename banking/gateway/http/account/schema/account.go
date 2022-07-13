package schema

import (
	"time"

	"github.com/thalissonfelipe/banking/banking/domain/entities"
	"github.com/thalissonfelipe/banking/banking/gateway/http/rest"
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

func MapToListAccountsResponse(accs []entities.Account) ListAccountsResponse {
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
	if r.Name == "" {
		return rest.ErrMissingNameParameter
	}

	if r.CPF == "" {
		return rest.ErrMissingCPFParameter
	}

	if r.Secret == "" {
		return rest.ErrMissingSecretParameter
	}

	return nil
}

type CreateAccountResponse struct {
	Name      string `json:"name"`
	CPF       string `json:"cpf"`
	Balance   int    `json:"balance"`
	CreatedAt string `json:"created_at"`
}

func MapToCreateAccountResponse(acc entities.Account) CreateAccountResponse {
	return CreateAccountResponse{
		Name:      acc.Name,
		CPF:       acc.CPF.String(),
		Balance:   acc.Balance,
		CreatedAt: acc.CreatedAt.UTC().Format(time.RFC3339),
	}
}
