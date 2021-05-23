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
	ErrInsufficientFunds              error = errors.New("insufficient funds")
	ErrAccountDestinationDoesNotExist error = errors.New("account destination does not exist")
)

func NewTransfer(accOriginID, accDestID vos.ID, amount int) Transfer {
	return Transfer{
		ID:                   vos.NewID(),
		AccountOriginID:      accOriginID,
		AccountDestinationID: accDestID,
		Amount:               amount,
		CreatedAt:            time.Now(),
	}
}
