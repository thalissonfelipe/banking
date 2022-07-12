package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"github.com/thalissonfelipe/banking/banking/domain/account/usecase"
	"github.com/thalissonfelipe/banking/banking/domain/entities"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
	"github.com/thalissonfelipe/banking/banking/tests/mocks"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
)

func TestAuthenticate(t *testing.T) {
	cpf := testdata.GetValidCPF()
	secret := testdata.GetValidSecret()

	testCases := []struct {
		name        string
		repo        *mocks.AccountRepositoryMock
		enc         *mocks.HashMock
		input       AuthenticateInput
		token       string // TODO: validate created token
		expectedErr error
	}{
		{
			name:        "should return an error if account does not exist",
			repo:        &mocks.AccountRepositoryMock{},
			enc:         &mocks.HashMock{},
			input:       AuthenticateInput{CPF: cpf.String(), Secret: secret.String()},
			expectedErr: ErrInvalidCredentials,
		},
		{
			name:        "should return an error if cpf provided is invalid",
			repo:        &mocks.AccountRepositoryMock{},
			enc:         &mocks.HashMock{},
			input:       AuthenticateInput{CPF: "123.456.789-00", Secret: secret.String()},
			expectedErr: vos.ErrInvalidCPF,
		},
		{
			name: "should return an error if secret does not match",
			repo: &mocks.AccountRepositoryMock{
				Accounts: []entities.Account{
					entities.NewAccount("Pedro", cpf, secret),
				},
			},
			enc:         &mocks.HashMock{Err: bcrypt.ErrMismatchedHashAndPassword},
			input:       AuthenticateInput{CPF: cpf.String(), Secret: secret.String()},
			expectedErr: ErrInvalidCredentials,
		},
		{
			name: "should return no error if authenticated succeeds",
			repo: &mocks.AccountRepositoryMock{
				Accounts: []entities.Account{
					entities.NewAccount("Pedro", cpf, secret),
				},
			},
			enc:         &mocks.HashMock{},
			input:       AuthenticateInput{CPF: cpf.String(), Secret: secret.String()},
			expectedErr: nil,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			accUsecase := usecase.NewAccountUsecase(tt.repo, nil)
			service := NewAuth(accUsecase, tt.enc)

			_, err := service.Autheticate(ctx, tt.input)

			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}
