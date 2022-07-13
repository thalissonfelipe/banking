package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/thalissonfelipe/banking/banking/domain/account"
	"github.com/thalissonfelipe/banking/banking/domain/account/usecase"
	"github.com/thalissonfelipe/banking/banking/domain/encrypter"
	"github.com/thalissonfelipe/banking/banking/domain/entities"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
)

func TestAuth_Authenticate(t *testing.T) {
	cpf := testdata.GetValidCPF()
	secret := testdata.GetValidSecret()

	acc, err := entities.NewAccount("name", cpf.String(), secret.String())
	require.NoError(t, err)

	testCases := []struct {
		name    string
		repo    account.Repository
		enc     encrypter.Encrypter
		cpf     string
		secret  string
		wantErr error
	}{
		{
			name: "should return a token if authentication succeeds",
			repo: &RepositoryMock{
				GetAccountByCPFFunc: func(context.Context, vos.CPF) (entities.Account, error) {
					return acc, nil
				},
			},
			enc: &EncrypterMock{
				CompareHashAndSecretFunc: func(_, _ []byte) error {
					return nil
				},
			},
			cpf:     cpf.String(),
			secret:  secret.String(),
			wantErr: nil,
		},
		{
			name: "should return an error if account does not exist",
			repo: &RepositoryMock{
				GetAccountByCPFFunc: func(context.Context, vos.CPF) (entities.Account, error) {
					return entities.Account{}, entities.ErrAccountDoesNotExist
				},
			},
			enc:     &EncrypterMock{},
			cpf:     cpf.String(),
			secret:  secret.String(),
			wantErr: ErrInvalidCredentials,
		},
		{
			name:    "should return an error if cpf provided is invalid",
			repo:    &RepositoryMock{},
			enc:     &EncrypterMock{},
			cpf:     "123.456.789-00",
			secret:  secret.String(),
			wantErr: vos.ErrInvalidCPF,
		},
		{
			name: "should return an error if secret does not match",
			repo: &RepositoryMock{
				GetAccountByCPFFunc: func(context.Context, vos.CPF) (entities.Account, error) {
					return acc, nil
				},
			},
			enc: &EncrypterMock{
				CompareHashAndSecretFunc: func(_, _ []byte) error {
					return bcrypt.ErrMismatchedHashAndPassword
				},
			},
			cpf:     cpf.String(),
			secret:  secret.String(),
			wantErr: ErrInvalidCredentials,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			accUsecase := usecase.NewAccountUsecase(tt.repo, tt.enc)
			service := NewAuth(accUsecase, tt.enc)

			_, err := service.Autheticate(context.Background(), tt.cpf, tt.secret)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}