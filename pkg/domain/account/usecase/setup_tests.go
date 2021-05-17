package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

type StubRepository struct {
	accounts []entities.Account
	err      error
}

func (s StubRepository) GetAccounts(ctx context.Context) ([]entities.Account, error) {
	if s.err != nil {
		return nil, entities.ErrInternalError
	}
	return s.accounts, nil
}

func (s StubRepository) GetBalanceByID(ctx context.Context, id string) (int, error) {
	for _, account := range s.accounts {
		if account.ID == id {
			return account.Balance, nil
		}
	}
	return 0, entities.ErrAccountDoesNotExist
}

func (s *StubRepository) PostAccount(ctx context.Context, account entities.Account) error {
	if s.err != nil {
		return entities.ErrInternalError
	}
	s.accounts = append(s.accounts, account)
	return nil
}

func (s StubRepository) GetAccountByCPF(ctx context.Context, cpf string) (*entities.Account, error) {
	for _, acc := range s.accounts {
		if acc.CPF == cpf {
			return &acc, s.err
		}
	}
	return nil, s.err
}

func (s StubRepository) GetAccountByID(ctx context.Context, id string) (*entities.Account, error) {
	for _, acc := range s.accounts {
		if acc.ID == id {
			return &acc, s.err
		}
	}
	return nil, s.err
}

type StubHash struct {
	err error
}

func (s StubHash) Hash(secret string) ([]byte, error) {
	if s.err != nil {
		return nil, entities.ErrInternalError
	}

	return []byte(generateRandomSecret(len(secret))), nil
}

func generateRandomSecret(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
