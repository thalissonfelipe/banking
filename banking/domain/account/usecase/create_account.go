package usecase

import (
	"context"
	"fmt"

	"github.com/thalissonfelipe/banking/banking/domain/account"
	"github.com/thalissonfelipe/banking/banking/domain/entities"
)

func (a Account) CreateAccount(ctx context.Context, input account.CreateAccountInput) (*entities.Account, error) {
	err := input.Secret.Hash(a.encrypter)
	if err != nil {
		return nil, fmt.Errorf("hashing secret: %w", err)
	}

	acc := entities.NewAccount(input.Name, input.CPF, input.Secret)

	err = a.repository.CreateAccount(ctx, &acc)
	if err != nil {
		return nil, fmt.Errorf("creating account: %w", err)
	}

	return &acc, nil
}
