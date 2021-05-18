package mocks

import (
	"context"

	"github.com/thalissonfelipe/banking/pkg/domain/account"
	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

type StubAccountUseCase struct {
	Accounts []entities.Account
	Err      error
}

func (s StubAccountUseCase) GetAccountBalanceByID(ctx context.Context, accountID string) (int, error) {
	for _, acc := range s.Accounts {
		if acc.ID == accountID {
			return acc.Balance, nil
		}
	}
	return 0, nil
}

func (s StubAccountUseCase) ListAccounts(ctx context.Context) ([]entities.Account, error) {
	return nil, nil
}

func (s StubAccountUseCase) CreateAccount(ctx context.Context, input account.CreateAccountInput) (*entities.Account, error) {
	return nil, nil
}

func (s StubAccountUseCase) GetAccountByID(ctx context.Context, accountID string) (*entities.Account, error) {
	for _, acc := range s.Accounts {
		if acc.ID == accountID {
			return &acc, nil
		}
	}
	return nil, nil
}
