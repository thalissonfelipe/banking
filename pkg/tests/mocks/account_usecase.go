package mocks

import (
	"context"

	"github.com/thalissonfelipe/banking/pkg/domain/account"
	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
)

var _ account.Usecase = (*AccountUsecaseMock)(nil)

type AccountUsecaseMock struct {
	Accounts []entities.Account
	Err      error
}

func (s AccountUsecaseMock) GetAccountBalanceByID(ctx context.Context, accountID vos.AccountID) (int, error) {
	if s.Err != nil {
		return 0, entities.ErrInternalError
	}

	for _, acc := range s.Accounts {
		if acc.ID == accountID {
			return acc.Balance, nil
		}
	}

	return 0, entities.ErrAccountDoesNotExist
}

func (s AccountUsecaseMock) ListAccounts(ctx context.Context) ([]entities.Account, error) {
	if s.Err != nil {
		return nil, entities.ErrInternalError
	}

	return s.Accounts, nil
}

func (s AccountUsecaseMock) CreateAccount(
	ctx context.Context, input account.CreateAccountInput) (*entities.Account, error) {
	if s.Err != nil {
		return nil, entities.ErrInternalError
	}

	for _, acc := range s.Accounts {
		if acc.CPF == input.CPF {
			return nil, entities.ErrAccountAlreadyExists
		}
	}

	acc := entities.NewAccount(
		input.Name,
		input.CPF,
		input.Secret,
	)

	s.Accounts = append(s.Accounts, acc)

	return &acc, nil
}

func (s AccountUsecaseMock) GetAccountByID(ctx context.Context, accountID vos.AccountID) (entities.Account, error) {
	if s.Err != nil {
		return entities.Account{}, entities.ErrInternalError
	}

	for _, acc := range s.Accounts {
		if acc.ID == accountID {
			return acc, nil
		}
	}

	return entities.Account{}, entities.ErrAccountDoesNotExist
}

func (s AccountUsecaseMock) GetAccountByCPF(ctx context.Context, cpf vos.CPF) (entities.Account, error) {
	if s.Err != nil {
		return entities.Account{}, entities.ErrInternalError
	}

	for _, acc := range s.Accounts {
		if acc.CPF == cpf {
			return acc, nil
		}
	}

	return entities.Account{}, entities.ErrAccountDoesNotExist
}
