package mocks

import (
	"context"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
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

func (s StubAccountRepository) GetBalanceByID(ctx context.Context, id vos.ID) (int, error) {
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

func (s *StubAccountRepository) CreateAccount(ctx context.Context, account *entities.Account) error {
	if s.Err != nil {
		return entities.ErrInternalError
	}

	for _, acc := range s.Accounts {
		if acc.CPF == account.CPF {
			return entities.ErrAccountAlreadyExists
		}
	}

	s.Accounts = append(s.Accounts, *account)

	return nil
}

func (s StubAccountRepository) GetAccountByCPF(ctx context.Context, cpf vos.CPF) (*entities.Account, error) {
	if s.Err != nil {
		return nil, s.Err
	}
	for _, acc := range s.Accounts {
		if acc.CPF == cpf {
			return &acc, nil
		}
	}
	return nil, entities.ErrAccountDoesNotExist
}

func (s StubAccountRepository) GetAccountByID(ctx context.Context, id vos.ID) (*entities.Account, error) {
	if s.Err != nil {
		return nil, s.Err
	}
	for _, acc := range s.Accounts {
		if acc.ID == id {
			return &acc, nil
		}
	}
	return nil, entities.ErrAccountDoesNotExist
}
