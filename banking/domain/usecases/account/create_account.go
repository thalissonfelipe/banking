package account

import (
	"context"
	"fmt"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
)

func (u Usecase) CreateAccount(ctx context.Context, acc *entity.Account) error {
	err := acc.Secret.Hash(u.encrypter)
	if err != nil {
		return fmt.Errorf("hashing secret: %w", err)
	}

	err = u.repository.CreateAccount(ctx, acc)
	if err != nil {
		return fmt.Errorf("creating account: %w", err)
	}

	return nil
}
