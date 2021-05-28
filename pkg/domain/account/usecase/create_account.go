package usecase

import (
	"context"
	"errors"

	"github.com/thalissonfelipe/banking/pkg/domain/account"
	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
)

func (a Account) CreateAccount(ctx context.Context, input account.CreateAccountInput) (*entities.Account, error) {
	_, err := a.repository.GetAccountByCPF(ctx, input.CPF)
	if errors.Is(err, nil) {
		return nil, entities.ErrAccountAlreadyExists
	}
	if !errors.Is(err, entities.ErrAccountDoesNotExist) {
		return nil, entities.ErrInternalError
	}

	hashedSecret, err := a.encrypter.Hash(input.Secret.String())
	if err != nil {
		return nil, entities.ErrInternalError
	}

	acc := entities.NewAccount(input.Name, input.CPF, vos.NewSecret(string(hashedSecret)))
	err = a.repository.PostAccount(ctx, &acc)
	if err != nil {
		return nil, entities.ErrInternalError
	}

	return &acc, nil
}
