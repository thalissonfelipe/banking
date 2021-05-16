package entities

import (
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

func NewAccountID() string {
	return uuid.New().String()
}

func NewAccount(name, secret, cpf string) Account {
	return Account{
		ID:        NewAccountID(),
		Name:      name,
		CPF:       cpf,
		Secret:    secret,
		Balance:   DefaultBalance,
		CreatedAt: time.Now(),
	}
}
