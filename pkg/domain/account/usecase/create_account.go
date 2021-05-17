package usecase

import (
	"context"

	"github.com/thalissonfelipe/banking/pkg/domain/account"
	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

func (a Account) CreateAccount(ctx context.Context, input account.CreateAccountInput) (*entities.Account, error) {
	acc := entities.NewAccount(input.Name, input.Secret, input.CPF)
	err := a.repository.PostAccount(ctx, acc)
	if err != nil {
		return nil, err
	}
	return &acc, err
}
