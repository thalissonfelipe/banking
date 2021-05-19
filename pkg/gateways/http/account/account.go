package account

import (
	"time"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
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

func convertAccountToAccountResponse(account entities.Account) accountResponse {
	return accountResponse{
		Name:      account.Name,
		CPF:       account.CPF.String(),
		Balance:   account.Balance,
		CreatedAt: formatTime(account.CreatedAt),
	}
}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05.000000")
}
