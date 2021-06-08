package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/thalissonfelipe/banking/pkg/domain/vos"
)

type Transfer struct {
	ID                   uuid.UUID
	AccountOriginID      vos.AccountID
	AccountDestinationID vos.AccountID
	Amount               int
	CreatedAt            time.Time
}

var (
	// ErrInsufficientFunds occurs when an account does not have sufficient funds.
	ErrInsufficientFunds = errors.New("insufficient funds")
	// ErrAccountDestinationDoesNotExist ocurrs when the account destination does not exist.
	ErrAccountDestinationDoesNotExist = errors.New("account destination does not exist")
)

func NewTransfer(accOriginID, accDestID vos.AccountID, amount int) Transfer {
	return Transfer{
		ID:                   uuid.New(),
		AccountOriginID:      accOriginID,
		AccountDestinationID: accDestID,
		Amount:               amount,
	}
}
