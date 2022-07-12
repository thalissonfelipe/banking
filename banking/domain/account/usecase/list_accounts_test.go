package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thalissonfelipe/banking/banking/domain/entities"
	"github.com/thalissonfelipe/banking/banking/tests/mocks"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
)

func TestUsecase_ListAccounts(t *testing.T) {
	acc := entities.NewAccount("Piter", testdata.GetValidCPF(), testdata.GetValidSecret())

	testCases := []struct {
		name      string
		repoSetup *mocks.AccountRepositoryMock
		want      []entities.Account
		wantErr   bool
	}{
		{
			name: "should return a list of accounts",
			repoSetup: &mocks.AccountRepositoryMock{
				Accounts: []entities.Account{acc},
			},
			want:    []entities.Account{acc},
			wantErr: false,
		},
		{
			name:      "should return an error if something went wrong on repository",
			repoSetup: &mocks.AccountRepositoryMock{Err: assert.AnError},
			want:      nil,
			wantErr:   true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			usecase := NewAccountUsecase(tt.repoSetup, nil)

			accounts, err := usecase.ListAccounts(context.Background())
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, accounts)
		})
	}
}
