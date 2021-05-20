package usecase

import (
	"context"
	"errors"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

func (a Account) GetAccountByCPF(ctx context.Context, cpf string) (*entities.Account, error) {
	acc, err := a.repository.GetAccountByCPF(ctx, cpf)
	if err != nil {
		if !errors.Is(err, entities.ErrAccountDoesNotExist) {
			return nil, entities.ErrInternalError
		}
		return nil, err
	}

	return acc, nil
}
