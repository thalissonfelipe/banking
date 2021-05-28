package mocks

import (
	"context"

	"github.com/thalissonfelipe/banking/pkg/domain/account"
	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
)

type StubAccountUsecase struct {
	Accounts []entities.Account
	Err      error
}

func (s StubAccountUsecase) GetAccountBalanceByID(ctx context.Context, accountID vos.ID) (int, error) {
	if s.Err != nil {
		return 0, entities.ErrInternalError
	}
	for _, acc := range s.Accounts {
		if acc.ID == accountID {
			return acc.Balance, nil
		}
	}

	return 0, nil
}

func (s StubAccountUsecase) ListAccounts(ctx context.Context) ([]entities.Account, error) {
	return nil, nil
}

func (s StubAccountUsecase) CreateAccount(ctx context.Context, input account.CreateAccountInput) (*entities.Account, error) {
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

func (s StubAccountUsecase) GetAccountByID(ctx context.Context, accountID vos.ID) (*entities.Account, error) {
	if s.Err != nil {
		return nil, entities.ErrInternalError
	}
	for _, acc := range s.Accounts {
		if acc.ID == accountID {
			return &acc, nil
		}
	}

	return nil, entities.ErrAccountDoesNotExist
}

func (s StubAccountUsecase) GetAccountByCPF(ctx context.Context, cpf vos.CPF) (*entities.Account, error) {
	if s.Err != nil {
		return nil, entities.ErrInternalError
	}
	for _, acc := range s.Accounts {
		if acc.CPF == cpf {
			return &acc, nil
		}
	}

	return nil, entities.ErrAccountDoesNotExist
}
