package rest

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

var (
	ErrMissingParameter = errors.New("missing parameter")
	ErrInvalidUUID      = errors.New("invalid uuid")
	ErrSameAccounts     = errors.New("account origin id cannot be equal to destination id")
	ErrUnauthorized     = errors.New("unauthorized")
)

// TODO: add tests.
type ValidationError struct {
	Location string
	Err      error
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", v.Location, v.Err.Error())
}

func (v ValidationError) Unwrap() error {
	return v.Err
}

// TODO: add tests.
type ValidationErrors []error

func (v ValidationErrors) Error() string {
	var (
		builder strings.Builder
		once    sync.Once
		sep     string
	)

	for _, err := range v {
		builder.WriteString(sep)
		builder.WriteString(err.Error())
		once.Do(func() {
			sep = "; "
		})
	}

	return builder.String()
}
