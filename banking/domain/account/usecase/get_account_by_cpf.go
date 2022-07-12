package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/thalissonfelipe/banking/banking/domain/entities"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
)

func (a Account) GetAccountByCPF(ctx context.Context, cpf vos.CPF) (entities.Account, error) {
	acc, err := a.repository.GetAccountByCPF(ctx, cpf)
	if err == nil {
		return acc, nil
	}

	if errors.Is(err, entities.ErrAccountDoesNotExist) {
		return entities.Account{}, fmt.Errorf("account does not exist: %w", err)
	}

	return entities.Account{}, entities.ErrInternalError
}
