package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

func TestListTransfers(t *testing.T) {
	ctx := context.Background()

	t.Run("should return a list of transfers", func(t *testing.T) {
		accountOriginID := entities.NewAccountID()
		transfer := entities.NewTransfer(
			accountOriginID,
			entities.NewAccountID(),
			100,
		)
		repo := StubRepository{transfers: []entities.Transfer{transfer}}
		usecase := NewTransfer(repo)
		expected := []entities.Transfer{transfer}
		result, err := usecase.ListTransfers(ctx, accountOriginID)

		assert.Nil(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("should return an error if something went wrong on repository", func(t *testing.T) {
		repo := StubRepository{transfers: nil, err: errors.New("failed to fetch transfers")}
		usecase := NewTransfer(repo)
		result, err := usecase.ListTransfers(ctx, entities.NewAccountID())

		assert.Nil(t, result)
		assert.NotNil(t, err)
	})
}
