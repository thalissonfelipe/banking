package entities

import (
	"errors"
	"time"

	"github.com/thalissonfelipe/banking/pkg/domain/vos"
)

type Transfer struct {
	ID                   vos.ID
	AccountOriginID      vos.ID
	AccountDestinationID vos.ID
	Amount               int
	CreatedAt            time.Time
}

var (
	// ErrInsufficientFunds occurs when an account does not have sufficient funds.
	ErrInsufficientFunds = errors.New("insufficient funds")
	// ErrAccountDestinationDoesNotExist ocurrs when the account destination does not exist.
	ErrAccountDestinationDoesNotExist = errors.New("account destination does not exist")
)

func NewTransfer(accOriginID, accDestID vos.ID, amount int) Transfer {
	return Transfer{
		ID:                   vos.NewID(),
		AccountOriginID:      accOriginID,
		AccountDestinationID: accDestID,
		Amount:               amount,
	}
}
