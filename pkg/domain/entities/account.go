package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
)

type Account struct {
	ID        string
	Name      string
	CPF       vos.CPF
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
	ErrInvalidCPF           error = errors.New("invalid cpf")
)

func NewAccountID() string {
	return uuid.New().String()
}

func NewAccount(name string, cpf vos.CPF, secret string) Account {
	return Account{
		ID:        NewAccountID(),
		Name:      name,
		CPF:       cpf,
		Secret:    secret,
		Balance:   DefaultBalance,
		CreatedAt: time.Now(),
	}
}
