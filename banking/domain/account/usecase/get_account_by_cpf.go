package usecase

import (
	"context"
	"fmt"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
)

func (a Account) GetAccountByCPF(ctx context.Context, cpf vos.CPF) (entity.Account, error) {
	acc, err := a.repository.GetAccountByCPF(ctx, cpf)
	if err != nil {
		return entity.Account{}, fmt.Errorf("getting account by cpf: %w", err)
	}

	return acc, nil
}
