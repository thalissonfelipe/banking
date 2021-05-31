package usecase

import (
	"context"
	"errors"

	"github.com/thalissonfelipe/banking/pkg/domain/account"
	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

func (a Account) CreateAccount(ctx context.Context, input account.CreateAccountInput) (*entities.Account, error) {
	err := input.Secret.Hash(a.encrypter)
	if err != nil {
		return nil, entities.ErrInternalError
	}

	acc := entities.NewAccount(input.Name, input.CPF, input.Secret)

	err = a.repository.CreateAccount(ctx, &acc)
	if err != nil {
		if errors.Is(err, entities.ErrAccountAlreadyExists) {
			return nil, err
		}
		return nil, entities.ErrInternalError
	}

	return &acc, nil
}
