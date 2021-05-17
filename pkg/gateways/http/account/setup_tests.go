package account

import (
	"context"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

type StubAccountRepository struct {
	accounts []entities.Account
	err      error
}

func (s StubAccountRepository) GetAccounts(ctx context.Context) ([]entities.Account, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.accounts, nil
}

func (s StubAccountRepository) GetBalanceByID(ctx context.Context, id string) (int, error) {
	for _, account := range s.accounts {
		if account.ID == id {
			return account.Balance, nil
		}
	}
	return 0, entities.ErrAccountDoesNotExist
}

func (s *StubAccountRepository) PostAccount(ctx context.Context, account entities.Account) error {
	if s.err != nil {
		return entities.ErrInternalError
	}
	s.accounts = append(s.accounts, account)
	return nil
}

func (s StubAccountRepository) GetAccountByCPF(ctx context.Context, cpf string) (*entities.Account, error) {
	for _, acc := range s.accounts {
		if acc.CPF == cpf {
			return &acc, s.err
		}
	}
	return nil, s.err
}

func (s StubAccountRepository) GetAccountByID(ctx context.Context, id string) (*entities.Account, error) {
	for _, acc := range s.accounts {
		if acc.ID == id {
			return &acc, s.err
		}
	}
	return nil, s.err
}
