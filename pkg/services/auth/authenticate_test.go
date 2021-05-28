package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thalissonfelipe/banking/pkg/domain/account/usecase"
	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
	"github.com/thalissonfelipe/banking/pkg/tests/mocks"
	"github.com/thalissonfelipe/banking/pkg/tests/testdata"
)

func TestAuthenticate(t *testing.T) {
	cpf := testdata.GetValidCPF()

	testCases := []struct {
		name  string
		repo  *mocks.StubAccountRepository
		enc   *mocks.StubHash
		input AuthenticateInput
		token string // TODO: validate created token
		err   error
	}{
		{
			name:  "should return an error if account does not exist",
			repo:  &mocks.StubAccountRepository{},
			enc:   &mocks.StubHash{},
			input: AuthenticateInput{CPF: cpf.String(), Secret: "12345678"},
			err:   entities.ErrAccountDoesNotExist,
		},
		{
			name:  "should return an error if cpf provided is invalid",
			repo:  &mocks.StubAccountRepository{},
			enc:   &mocks.StubHash{},
			input: AuthenticateInput{CPF: "123.456.789-00", Secret: "12345678"},
			err:   entities.ErrInvalidCPF,
		},
		{
			name: "should return an error if usecase fails to fetch account",
			repo: &mocks.StubAccountRepository{
				Err: errors.New("usecase fails to fetch account"),
			},
			enc:   &mocks.StubHash{},
			input: AuthenticateInput{CPF: cpf.String(), Secret: "12345678"},
			err:   entities.ErrInternalError,
		},
		{
			name: "should return an error if secret does not match",
			repo: &mocks.StubAccountRepository{
				Accounts: []entities.Account{
					entities.NewAccount("Pedro", cpf, vos.NewSecret("12345678")),
				},
			},
			enc: &mocks.StubHash{
				Err: ErrSecretDoesNotMatch,
			},
			input: AuthenticateInput{CPF: cpf.String(), Secret: "12345678"},
			err:   ErrSecretDoesNotMatch,
		},
		{
			name: "should return nil if authenticated succeeds",
			repo: &mocks.StubAccountRepository{
				Accounts: []entities.Account{
					entities.NewAccount("Pedro", cpf, vos.NewSecret("12345678")),
				},
			},
			enc:   &mocks.StubHash{},
			input: AuthenticateInput{CPF: cpf.String(), Secret: "12345678"},
			err:   nil,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			accUsecase := usecase.NewAccountUsecase(tt.repo, nil)
			service := NewAuth(accUsecase, tt.enc)

			_, err := service.Autheticate(ctx, tt.input)

			assert.Equal(t, tt.err, err)
		})
	}
}
