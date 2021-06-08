package mocks

import (
	"context"

	"github.com/thalissonfelipe/banking/pkg/domain/account"
	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
)

var _ account.Repository = (*AccountRepositoryMock)(nil)

type AccountRepositoryMock struct {
	Accounts []entities.Account
	Err      error
}

func (s AccountRepositoryMock) GetAccounts(ctx context.Context) ([]entities.Account, error) {
	if s.Err != nil {
		return nil, entities.ErrInternalError
	}

	return s.Accounts, nil
}

func (s AccountRepositoryMock) GetBalanceByID(ctx context.Context, id vos.AccountID) (int, error) {
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

func (s *AccountRepositoryMock) CreateAccount(ctx context.Context, account *entities.Account) error {
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

func (s AccountRepositoryMock) GetAccountByCPF(ctx context.Context, cpf vos.CPF) (*entities.Account, error) {
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

func (s AccountRepositoryMock) GetAccountByID(ctx context.Context, id vos.AccountID) (*entities.Account, error) {
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
