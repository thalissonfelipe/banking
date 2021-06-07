package account

import (
	"errors"
	"time"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

var (
	errMissingNameParameter   = errors.New("missing name parameter")
	errMissingCPFParameter    = errors.New("missing cpf parameter")
	errMissingSecretParameter = errors.New("missing secret parameter")
)

type AccountListResponse struct {
	Id        string `json:"id"`
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

func (r CreateAccountInput) isValid() error {
	if r.Name == "" {
		return errMissingNameParameter
	}

	if r.CPF == "" {
		return errMissingCPFParameter
	}

	if r.Secret == "" {
		return errMissingSecretParameter
	}

	return nil
}

func convertAccountToAccountListResponse(account entities.Account) AccountListResponse {
	return AccountListResponse{
		Id:        account.ID.String(),
		Name:      account.Name,
		CPF:       account.CPF.String(),
		Balance:   account.Balance,
		CreatedAt: formatTime(account.CreatedAt),
	}
}

func convertAccountToCreateAccountResponse(account *entities.Account) CreateAccountResponse {
	return CreateAccountResponse{
		Name:    account.Name,
		CPF:     account.CPF.String(),
		Balance: account.Balance,
	}
}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05.000000")
}
