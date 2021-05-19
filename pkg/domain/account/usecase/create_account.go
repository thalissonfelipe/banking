package usecase

import (
	"context"

	"github.com/thalissonfelipe/banking/pkg/domain/account"
	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
)

func (a Account) CreateAccount(ctx context.Context, input account.CreateAccountInput) (*entities.Account, error) {
	secret := vos.NewSecret(input.Secret)
	if ok := secret.IsValid(); !ok {
		return nil, entities.ErrInvalidSecret
	}

	accExists, err := a.repository.GetAccountByCPF(ctx, input.CPF.String())
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

	acc := entities.NewAccount(input.Name, input.CPF, string(hashedSecret))
	err = a.repository.PostAccount(ctx, acc)
	if err != nil {
		return nil, entities.ErrInternalError
	}

	return &acc, nil
}
