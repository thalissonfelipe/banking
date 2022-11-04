package transfer

import (
	"context"
	"fmt"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
)

func (u Usecase) ListTransfers(ctx context.Context, accountID vos.AccountID) ([]entity.Transfer, error) {
	transfers, err := u.repository.ListTransfers(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("listing transfers: %w", err)
	}

	return transfers, nil
}
