package mocks

import (
	"context"

	"github.com/thalissonfelipe/banking/pkg/domain/account"
	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
)

type StubAccountUseCase struct {
	Accounts []entities.Account
	Err      error
}

func (s StubAccountUseCase) GetAccountBalanceByID(ctx context.Context, accountID vos.ID) (int, error) {
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

func (s StubAccountUseCase) ListAccounts(ctx context.Context) ([]entities.Account, error) {
	return nil, nil
}

func (s StubAccountUseCase) CreateAccount(ctx context.Context, input account.CreateAccountInput) (*entities.Account, error) {
	return nil, nil
}

func (s StubAccountUseCase) GetAccountByID(ctx context.Context, accountID vos.ID) (*entities.Account, error) {
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

func (s StubAccountUseCase) GetAccountByCPF(ctx context.Context, cpf string) (*entities.Account, error) {
	if s.Err != nil {
		return nil, entities.ErrInternalError
	}
	for _, acc := range s.Accounts {
		if acc.CPF.String() == cpf {
			return &acc, nil
		}
	}

	return nil, entities.ErrAccountDoesNotExist
}
