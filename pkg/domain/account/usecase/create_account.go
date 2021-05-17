package usecase

import (
	"context"

	"github.com/thalissonfelipe/banking/pkg/domain/account"
	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

func (a Account) CreateAccount(ctx context.Context, input account.CreateAccountInput) (*entities.Account, error) {
	accExists, err := a.repository.GetAccountByCPF(ctx, input.CPF)
	if err != nil {
		return nil, entities.ErrInternalError
	}
	if accExists != nil {
		return nil, entities.ErrAccountAlreadyExists
	}

	hashedSecret, err := a.encrypter.Hash(input.Secret)
	if err != nil {
		return nil, entities.ErrInternalError
	}

	input.Secret = string(hashedSecret)
	acc := entities.NewAccount(input.Name, input.Secret, input.CPF)
	err = a.repository.PostAccount(ctx, acc)
	if err != nil {
		return nil, entities.ErrInternalError
	}

	return &acc, nil
}
