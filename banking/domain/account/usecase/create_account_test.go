package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/banking/domain/account"
	"github.com/thalissonfelipe/banking/banking/domain/encrypter"
	"github.com/thalissonfelipe/banking/banking/domain/entities"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
)

func TestAccountUsecase_CreateAccount(t *testing.T) {
	testCases := []struct {
		name    string
		repo    account.Repository
		enc     encrypter.Encrypter
		wantErr error
	}{
		{
			name: "should create an account successfully",
			repo: &RepositoryMock{
				CreateAccountFunc: func(context.Context, *entities.Account) error {
					return nil
				},
			},
			enc: &EncrypterMock{
				HashFunc: func(string) ([]byte, error) {
					return nil, nil
				},
			},
			wantErr: nil,
		},
		{
			name: "should return an error if encrypter fails to hash secret",
			repo: &RepositoryMock{},
			enc: &EncrypterMock{
				HashFunc: func(string) ([]byte, error) {
					return nil, assert.AnError
				},
			},
			wantErr: assert.AnError,
		},
		{
			name: "should return an error if account already exists",
			repo: &RepositoryMock{
				CreateAccountFunc: func(context.Context, *entities.Account) error {
					return entities.ErrAccountAlreadyExists
				},
			},
			enc: &EncrypterMock{
				HashFunc: func(string) ([]byte, error) {
					return nil, nil
				},
			},
			wantErr: entities.ErrAccountAlreadyExists,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			usecase := NewAccountUsecase(tt.repo, tt.enc)

			acc, err := entities.NewAccount("name", testdata.GetValidCPF().String(), testdata.GetValidSecret().String())
			require.NoError(t, err)

			err = usecase.CreateAccount(context.Background(), &acc)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
