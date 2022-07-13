package usecase

import (
	"context"
	"fmt"

	"github.com/thalissonfelipe/banking/banking/domain/entities"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
)

func (a Account) GetAccountByCPF(ctx context.Context, cpf vos.CPF) (entities.Account, error) {
	acc, err := a.repository.GetAccountByCPF(ctx, cpf)
	if err != nil {
		return entities.Account{}, fmt.Errorf("getting account by cpf: %w", err)
	}

	return acc, nil
}
