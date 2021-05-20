package usecase

import (
	"context"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

func (a Account) GetAccountByCPF(ctx context.Context, cpf string) (*entities.Account, error) {
	acc, err := a.repository.GetAccountByCPF(ctx, cpf)
	if err != nil {
		return nil, entities.ErrInternalError
	}
	if acc == nil {
		return nil, entities.ErrAccountDoesNotExist
	}

	return acc, nil
}
