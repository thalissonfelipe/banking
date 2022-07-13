package usecase

import (
	"context"
	"fmt"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
)

func (a Account) CreateAccount(ctx context.Context, acc *entity.Account) error {
	err := acc.Secret.Hash(a.encrypter)
	if err != nil {
		return fmt.Errorf("hashing secret: %w", err)
	}

	err = a.repository.CreateAccount(ctx, acc)
	if err != nil {
		return fmt.Errorf("creating account: %w", err)
	}

	return nil
}
