package entities

import (
	"errors"
	"time"

	"github.com/thalissonfelipe/banking/pkg/domain/vos"
)

type Account struct {
	ID        vos.ID
	Name      string
	CPF       vos.CPF
	Secret    vos.Secret
	Balance   int
	CreatedAt time.Time
}

const DefaultBalance = 0

var (
	ErrAccountDoesNotExist  error = errors.New("account does not exist")
	ErrAccountAlreadyExists error = errors.New("account already exists")
	ErrInternalError        error = errors.New("internal server error")
)

func NewAccount(name string, cpf vos.CPF, secret vos.Secret) Account {
	return Account{
		ID:      vos.NewID(),
		Name:    name,
		CPF:     cpf,
		Secret:  secret,
		Balance: DefaultBalance,
	}
}
