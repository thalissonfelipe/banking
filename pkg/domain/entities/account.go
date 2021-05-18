package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID        string
	Name      string
	CPF       string
	Secret    string
	Balance   int
	CreatedAt time.Time
}

const DefaultBalance = 0

var (
	ErrAccountDoesNotExist  error = errors.New("account does not exist")
	ErrAccountAlreadyExists error = errors.New("account already exists")
	ErrInternalError        error = errors.New("internal server error")
	ErrInvalidSecret        error = errors.New("invalid secret")
)

func NewAccountID() string {
	return uuid.New().String()
}

func NewAccount(name, cpf, secret string) Account {
	return Account{
		ID:        NewAccountID(),
		Name:      name,
		CPF:       cpf,
		Secret:    secret,
		Balance:   DefaultBalance,
		CreatedAt: time.Now(),
	}
}
