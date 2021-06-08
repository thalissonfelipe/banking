package entities

import (
	"errors"
	"time"

	"github.com/thalissonfelipe/banking/pkg/domain/vos"
)

type Account struct {
	ID        vos.AccountID
	Name      string
	CPF       vos.CPF
	Secret    vos.Secret
	Balance   int
	CreatedAt time.Time
}

const DefaultBalance = 0

var (
	// ErrAccountDoesNotExist occurs when an account does not exist.
	ErrAccountDoesNotExist = errors.New("account does not exist")
	// ErrAccountAlreadyExists occurs when an account already exists.
	ErrAccountAlreadyExists = errors.New("account already exists")
	// ErrInternalError ocurrs when an unexpected error happens.
	ErrInternalError = errors.New("internal server error")
)

func NewAccount(name string, cpf vos.CPF, secret vos.Secret) Account {
	return Account{
		ID:      vos.NewAccountID(),
		Name:    name,
		CPF:     cpf,
		Secret:  secret,
		Balance: DefaultBalance,
	}
}
