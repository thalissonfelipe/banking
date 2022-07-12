package vos

import (
	"log"

	"github.com/google/uuid"
)

type AccountID uuid.UUID

func (i AccountID) String() string {
	return uuid.UUID(i).String()
}

func NewAccountID() AccountID {
	return AccountID(uuid.New())
}

func ConvertStringToAccountID(id string) AccountID {
	uuID, err := uuid.Parse(id)
	if err != nil {
		log.Printf("unable to parse uuid: %s", err.Error())
	}

	return AccountID(uuID)
}
