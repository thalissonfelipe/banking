package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Transfer struct {
	ID                   string
	AccountOriginID      string
	AccountDestinationID string
	Amount               int
	CreatedAt            time.Time
}

var (
	ErrInsufficientFunds              error = errors.New("insufficient funds")
	ErrAccountDestinationDoesNotExist error = errors.New("account destination does not exist")
)

func NewTransferID() string {
	return uuid.New().String()
}

func NewTransfer(accOriginID, accDestID string, amount int) Transfer {
	return Transfer{
		ID:                   NewTransferID(),
		AccountOriginID:      accOriginID,
		AccountDestinationID: accDestID,
		Amount:               amount,
		CreatedAt:            time.Now(),
	}
}
