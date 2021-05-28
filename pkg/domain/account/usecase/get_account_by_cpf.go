package usecase

import (
	"context"
	"errors"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
)

func (a Account) GetAccountByCPF(ctx context.Context, cpf vos.CPF) (*entities.Account, error) {
	acc, err := a.repository.GetAccountByCPF(ctx, cpf)
	if err != nil {
		if !errors.Is(err, entities.ErrAccountDoesNotExist) {
			return nil, entities.ErrInternalError
		}
		return nil, err
	}

	return acc, nil
}
