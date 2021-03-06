package schema

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/gateway/http/rest"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
)

func TestSchema_MapToListAccountsResponse(t *testing.T) {
	t.Parallel()

	acc, err := entity.NewAccount("name", testdata.CPF().String(), testdata.Secret().String())
	require.NoError(t, err)

	acc.CreatedAt = time.Now()

	tests := []struct {
		name     string
		accounts []entity.Account
		want     ListAccountsResponse
	}{
		{
			name:     "empty list of accounts",
			accounts: nil,
			want:     ListAccountsResponse{Accounts: []Account{}},
		},
		{
			name:     "map accounts successfully",
			accounts: []entity.Account{acc, acc, acc, acc},
			want: ListAccountsResponse{Accounts: []Account{
				{
					ID:        acc.ID.String(),
					Name:      acc.Name,
					CPF:       acc.CPF.String(),
					Balance:   acc.Balance,
					CreatedAt: acc.CreatedAt.UTC().Format(time.RFC3339),
				},
				{
					ID:        acc.ID.String(),
					Name:      acc.Name,
					CPF:       acc.CPF.String(),
					Balance:   acc.Balance,
					CreatedAt: acc.CreatedAt.UTC().Format(time.RFC3339),
				},
				{
					ID:        acc.ID.String(),
					Name:      acc.Name,
					CPF:       acc.CPF.String(),
					Balance:   acc.Balance,
					CreatedAt: acc.CreatedAt.UTC().Format(time.RFC3339),
				},
				{
					ID:        acc.ID.String(),
					Name:      acc.Name,
					CPF:       acc.CPF.String(),
					Balance:   acc.Balance,
					CreatedAt: acc.CreatedAt.UTC().Format(time.RFC3339),
				},
			}},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := MapToListAccountsResponse(tt.accounts)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSchema_MapToBalanceResponse(t *testing.T) {
	t.Parallel()

	const balance = 100

	got := MapToBalanceResponse(balance)
	assert.Equal(t, BalanceResponse{Balance: balance}, got)
}

func TestSchema_CreateAccountInput_IsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   CreateAccountInput
		wantErr rest.ValidationErrors
	}{
		{
			name: "should validate input without errors",
			input: CreateAccountInput{
				Name:   "name",
				CPF:    "cpf",
				Secret: "secret",
			},
			wantErr: nil,
		},
		{
			name: "should return err if name is blank",
			input: CreateAccountInput{
				CPF:    "cpf",
				Secret: "secret",
			},
			wantErr: rest.ValidationErrors{ErrMissingNameParameter},
		},
		{
			name: "should return err if cpf is blank",
			input: CreateAccountInput{
				Name:   "name",
				Secret: "secret",
			},
			wantErr: rest.ValidationErrors{ErrMissingCPFParameter},
		},
		{
			name: "should return err if secret is blank",
			input: CreateAccountInput{
				Name: "name",
				CPF:  "cpf",
			},
			wantErr: rest.ValidationErrors{ErrMissingSecretParameter},
		},
		{
			name:    "should return err if all fields are blank",
			input:   CreateAccountInput{},
			wantErr: rest.ValidationErrors{ErrMissingNameParameter, ErrMissingCPFParameter, ErrMissingSecretParameter},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.input.IsValid()
			if err != nil {
				var errs rest.ValidationErrors
				require.True(t, errors.As(err, &errs))

				assert.Len(t, errs, len(tt.wantErr))

				for i, e := range errs {
					var verr rest.ValidationError
					require.True(t, errors.As(e, &verr))

					assert.ErrorIs(t, verr, tt.wantErr[i])
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
