package vos

import (
	"github.com/google/uuid"
)

type AccountID uuid.UUID

func (i AccountID) String() string {
	return uuid.UUID(i).String()
}

func NewAccountID() AccountID {
	return AccountID(uuid.New())
}
