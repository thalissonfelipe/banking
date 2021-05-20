package mocks

import (
	"context"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

type StubAccountRepository struct {
	Accounts []entities.Account
	Err      error
}

func (s StubAccountRepository) GetAccounts(ctx context.Context) ([]entities.Account, error) {
	if s.Err != nil {
		return nil, entities.ErrInternalError
	}
	return s.Accounts, nil
}

func (s StubAccountRepository) GetBalanceByID(ctx context.Context, id string) (int, error) {
	if s.Err != nil {
		return 0, s.Err
	}
	for _, account := range s.Accounts {
		if account.ID == id {
			return account.Balance, nil
		}
	}
	return 0, entities.ErrAccountDoesNotExist
}

func (s *StubAccountRepository) PostAccount(ctx context.Context, account entities.Account) error {
	if s.Err != nil {
		return entities.ErrInternalError
	}
	s.Accounts = append(s.Accounts, account)
	return nil
}

func (s StubAccountRepository) GetAccountByCPF(ctx context.Context, cpf string) (*entities.Account, error) {
	if s.Err != nil {
		return nil, entities.ErrInternalError
	}
	for _, acc := range s.Accounts {
		if acc.CPF.String() == cpf {
			return &acc, nil
		}
	}
	return nil, nil
}

func (s StubAccountRepository) GetAccountByID(ctx context.Context, id string) (*entities.Account, error) {
	if s.Err != nil {
		return nil, entities.ErrInternalError
	}
	for _, acc := range s.Accounts {
		if acc.ID == id {
			return &acc, nil
		}
	}
	return nil, nil
}
