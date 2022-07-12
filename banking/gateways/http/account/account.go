package account

import (
	"time"

	"github.com/thalissonfelipe/banking/banking/domain/entities"
	"github.com/thalissonfelipe/banking/banking/gateways/http/account/schemes"
)

func convertAccountToAccountListResponse(account entities.Account) schemes.AccountListResponse {
	return schemes.AccountListResponse{
		ID:        account.ID.String(),
		Name:      account.Name,
		CPF:       account.CPF.String(),
		Balance:   account.Balance,
		CreatedAt: formatTime(account.CreatedAt),
	}
}

func convertAccountToCreateAccountResponse(account *entities.Account) schemes.CreateAccountResponse {
	return schemes.CreateAccountResponse{
		Name:    account.Name,
		CPF:     account.CPF.String(),
		Balance: account.Balance,
	}
}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05.000000")
}
