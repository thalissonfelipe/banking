package usecase

import (
	"context"

	"github.com/thalissonfelipe/banking/pkg/domain/account"
	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

type StubRepository struct {
	transfers []entities.Transfer
	err       error
}

func (s StubRepository) GetTransfers(ctx context.Context, id string) ([]entities.Transfer, error) {
	if s.err != nil {
		return nil, entities.ErrInternalError
	}

	var transfers []entities.Transfer
	for _, tr := range s.transfers {
		if tr.AccountOriginID == id {
			transfers = append(transfers, tr)
		}
	}

	return transfers, nil
}

func (s *StubRepository) UpdateBalance(ctx context.Context, transfer entities.Transfer) error {
	s.transfers = append(s.transfers, transfer)
	return nil
}

type StubAccountUseCase struct {
	accounts []entities.Account
}

func (s StubAccountUseCase) GetAccountBalanceByID(ctx context.Context, accountID string) (int, error) {
	for _, acc := range s.accounts {
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
