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
	errInvalidJSON            = errors.New("invalid json")
	errInvalidCPF             = errors.New("invalid cpf")
	errInvalidSecret          = errors.New("invalid secret")
)

type accountResponse struct {
	Name      string
	CPF       string
	Balance   int
	CreatedAt string
}

type balanceResponse struct {
	Balance int
}

type requestBody struct {
	Name   string
	CPF    string
	Secret string
}

type createdAccountResponse struct {
	Name    string
	CPF     string
	Balance int
}

func (r requestBody) isValid() error {
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

func convertAccountToAccountResponse(account entities.Account) accountResponse {
	return accountResponse{
		Name:      account.Name,
		CPF:       account.CPF.String(),
		Balance:   account.Balance,
		CreatedAt: formatTime(account.CreatedAt),
	}
}

func convertAccountToCreatedAccountResponse(account *entities.Account) createdAccountResponse {
	return createdAccountResponse{
		Name:    account.Name,
		CPF:     account.CPF.String(),
		Balance: account.Balance,
	}
}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05.000000")
}
