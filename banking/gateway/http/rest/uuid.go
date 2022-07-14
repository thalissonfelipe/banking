package rest

import "github.com/google/uuid"

// TODO: add tests.
func ParseUUID(id, loc string) (uuid.UUID, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return uuid.UUID{}, ValidationError{Location: loc, Err: ErrInvalidUUID}
	}

	return parsedID, nil
}
