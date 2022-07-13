package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/thalissonfelipe/banking/banking/domain/vos"
)

type Transfer struct {
	ID                   uuid.UUID
	AccountOriginID      vos.AccountID
	AccountDestinationID vos.AccountID
	Amount               int
	CreatedAt            time.Time
}

// ErrInsufficientFunds occurs when an account does not have sufficient funds.
var ErrInsufficientFunds = errors.New("insufficient funds")

func NewTransfer(accOriginID, accDestID vos.AccountID, amount, accOriginBalance int) (Transfer, error) {
	if (accOriginBalance - amount) < 0 {
		return Transfer{}, ErrInsufficientFunds
	}

	return Transfer{
		ID:                   uuid.New(),
		AccountOriginID:      accOriginID,
		AccountDestinationID: accDestID,
		Amount:               amount,
	}, nil
}
