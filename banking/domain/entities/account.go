package entities

import (
	"errors"
	"time"

	"github.com/thalissonfelipe/banking/banking/domain/vos"
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
)

func NewAccount(name, cpfStr, secretStr string) (Account, error) {
	cpf, err := vos.NewCPF(cpfStr)
	if err != nil {
		return Account{}, err
	}

	secret, err := vos.NewSecret(secretStr)
	if err != nil {
		return Account{}, err
	}

	return Account{
		ID:      vos.NewAccountID(),
		Name:    name,
		CPF:     cpf,
		Secret:  secret,
		Balance: DefaultBalance,
	}, nil
}
