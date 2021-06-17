package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
)

func (a Account) GetAccountByCPF(ctx context.Context, cpf vos.CPF) (*entities.Account, error) {
	acc, err := a.repository.GetAccountByCPF(ctx, cpf)
	if err == nil {
		return acc, nil
	}

	if errors.Is(err, entities.ErrAccountDoesNotExist) {
		return nil, fmt.Errorf("account does not exist: %w", err)
	}

	return nil, entities.ErrInternalError
}
