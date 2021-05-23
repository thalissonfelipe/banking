package vos

import (
	"log"

	"github.com/google/uuid"
)

type ID uuid.UUID

func (i ID) String() string {
	return uuid.UUID(i).String()
}

func NewID() ID {
	return ID(uuid.New())
}

func ConvertStringToID(id string) ID {
	uuID, err := uuid.Parse(id)
	if err != nil {
		log.Printf("unable to parse uuid: %s", err.Error())
	}

	return ID(uuID)
}
