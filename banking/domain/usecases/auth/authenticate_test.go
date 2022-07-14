package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/banking/domain/encrypter"
	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/domain/usecases"
	"github.com/thalissonfelipe/banking/banking/domain/usecases/account"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
)

func TestAuth_Authenticate(t *testing.T) {
	cpf := testdata.GetValidCPF()
	secret := testdata.GetValidSecret()

	acc, err := entity.NewAccount("name", cpf.String(), secret.String())
	require.NoError(t, err)

	testCases := []struct {
		name      string
		repo      entity.AccountRepository
		enc       encrypter.Encrypter
		service   Service
		cpf       string
		secret    string
		wantToken string
		wantErr   error
	}{
		{
			name: "should return a token if authentication succeeds",
			repo: &RepositoryMock{
				GetAccountByCPFFunc: func(context.Context, vos.CPF) (entity.Account, error) {
					return acc, nil
				},
			},
			enc: &EncrypterMock{
				CompareHashAndSecretFunc: func(_, _ []byte) error {
					return nil
				},
			},
			service: &ServiceMock{
				NewTokenFunc: func(string) (string, error) {
					return "token", nil
				},
			},
			cpf:       cpf.String(),
			secret:    secret.String(),
			wantToken: "token",
			wantErr:   nil,
		},
		{
			name: "should return an error if account does not exist",
			repo: &RepositoryMock{
				GetAccountByCPFFunc: func(context.Context, vos.CPF) (entity.Account, error) {
					return entity.Account{}, entity.ErrAccountNotFound
				},
			},
			enc:       &EncrypterMock{},
			service:   &ServiceMock{},
			cpf:       cpf.String(),
			secret:    secret.String(),
			wantToken: "",
			wantErr:   usecases.ErrInvalidCredentials,
		},
		{
			name:      "should return an error if cpf provided is invalid",
			repo:      &RepositoryMock{},
			enc:       &EncrypterMock{},
			service:   &ServiceMock{},
			cpf:       "123.456.789-00",
			secret:    secret.String(),
			wantToken: "",
			wantErr:   vos.ErrInvalidCPF,
		},
		{
			name: "should return an error if secret does not match",
			repo: &RepositoryMock{
				GetAccountByCPFFunc: func(context.Context, vos.CPF) (entity.Account, error) {
					return acc, nil
				},
			},
			enc: &EncrypterMock{
				CompareHashAndSecretFunc: func(_, _ []byte) error {
					return usecases.ErrInvalidCredentials
				},
			},
			service:   &ServiceMock{},
			cpf:       cpf.String(),
			secret:    secret.String(),
			wantToken: "",
			wantErr:   usecases.ErrInvalidCredentials,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			accUsecase := account.NewAccountUsecase(tt.repo, tt.enc)
			service := NewAuth(accUsecase, tt.enc, tt.service)

			token, err := service.Autheticate(context.Background(), tt.cpf, tt.secret)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantToken, token)
		})
	}
}
