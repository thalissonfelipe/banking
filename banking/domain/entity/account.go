package entity

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/multierr"

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
	var errs error

	cpf, err := vos.NewCPF(cpfStr)
	if err != nil {
		errs = multierr.Append(errs, err)
	}

	secret, err := vos.NewSecret(secretStr)
	if err != nil {
		errs = multierr.Append(errs, err)
	}

	if errs != nil {
		return Account{}, fmt.Errorf("new account: %w", errs)
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
