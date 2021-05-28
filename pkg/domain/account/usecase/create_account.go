package usecase

import (
	"context"
	"errors"

	"github.com/thalissonfelipe/banking/pkg/domain/account"
	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
)

func (a Account) CreateAccount(ctx context.Context, input account.CreateAccountInput) (*entities.Account, error) {
	hashedSecret, err := a.encrypter.Hash(input.Secret.String())
	if err != nil {
		return nil, entities.ErrInternalError
	}

	acc := entities.NewAccount(input.Name, input.CPF, vos.NewSecret(string(hashedSecret)))

	err = a.repository.CreateAccount(ctx, &acc)
	if err != nil {
		if errors.Is(err, entities.ErrAccountAlreadyExists) {
			return nil, err
		}
		return nil, entities.ErrInternalError
	}

	return &acc, nil
}
