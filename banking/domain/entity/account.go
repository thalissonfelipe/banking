package entity

import (
	"context"
	"errors"
	"time"

	"github.com/thalissonfelipe/banking/banking/domain/vos"
)

var (
	// ErrAccountNotFound occurs when an account does not exist.
	ErrAccountNotFound = errors.New("account does not exist")
	// ErrAccountDestinationNotFound ocurrs when the account destination does not exist.
	ErrAccountDestinationNotFound = errors.New("account destination does not exist")
	// ErrAccountAlreadyExists occurs when an account already exists.
	ErrAccountAlreadyExists = errors.New("account already exists")
)

type Account struct {
	ID        vos.AccountID
	Name      string
	CPF       vos.CPF
	Secret    vos.Secret
	Balance   int
	CreatedAt time.Time
}

const defaultBalance = 100

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
		Balance: defaultBalance,
	}, nil
}

type AccountRepository interface {
	ListAccounts(context.Context) ([]Account, error)
	GetAccountBalanceByID(context.Context, vos.AccountID) (int, error)
	CreateAccount(context.Context, *Account) error
	GetAccountByCPF(context.Context, vos.CPF) (Account, error)
	GetAccountByID(context.Context, vos.AccountID) (Account, error)
}
