package schema

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/gateway/http/rest"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
)

func TestSchema_MapToListAccountsResponse(t *testing.T) {
	acc, err := entity.NewAccount("name", testdata.GetValidCPF().String(), testdata.GetValidSecret().String())
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
		t.Run(tt.name, func(t *testing.T) {
			got := MapToListAccountsResponse(tt.accounts)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSchema_MapToBalanceResponse(t *testing.T) {
	const balance = 100

	got := MapToBalanceResponse(balance)
	assert.Equal(t, BalanceResponse{Balance: balance}, got)
}

func TestSchema_CreateAccountInput_IsValid(t *testing.T) {
	tests := []struct {
		name    string
		input   CreateAccountInput
		wantErr error
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
			wantErr: rest.ErrMissingNameParameter,
		},
		{
			name: "should return err if cpf is blank",
			input: CreateAccountInput{
				Name:   "name",
				Secret: "secret",
			},
			wantErr: rest.ErrMissingCPFParameter,
		},
		{
			name: "should return err if secret is blank",
			input: CreateAccountInput{
				Name: "name",
				CPF:  "cpf",
			},
			wantErr: rest.ErrMissingSecretParameter,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.IsValid()
			assert.ErrorIs(t, got, tt.wantErr)
		})
	}
}