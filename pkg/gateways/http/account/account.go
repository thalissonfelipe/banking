package account

import (
	"time"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

type AccountResponse struct {
	Name      string
	CPF       string
	Balance   int
	CreatedAt string
}

type BalanceResponse struct {
	Balance int
}

func convertAccountToAccountResponse(account entities.Account) AccountResponse {
	return AccountResponse{
		Name:      account.Name,
		CPF:       account.CPF,
		Balance:   account.Balance,
		CreatedAt: formatTime(account.CreatedAt),
	}
}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05.000000")
}
