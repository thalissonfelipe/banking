package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thalissonfelipe/banking/banking/domain/vos"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
)

func TestNewAccount(t *testing.T) {
	t.Parallel()

	cpf := testdata.CPF()
	secret := testdata.Secret()

	tests := []struct {
		name     string
		cpf      string
		secret   string
		want     Account
		wantErrs []error
	}{
		{
			name:   "should create an account successfully",
			cpf:    cpf.String(),
			secret: secret.String(),
			want: Account{
				Name:    "name",
				CPF:     cpf,
				Secret:  secret,
				Balance: 100,
			},
			wantErrs: nil,
		},
		{
			name:     "should return an error if cpf is invalid",
			cpf:      "12345678900",
			secret:   secret.String(),
			want:     Account{},
			wantErrs: []error{vos.ErrInvalidCPF},
		},
		{
			name:     "should return an error if secret is invalid",
			cpf:      cpf.String(),
			secret:   "invalid",
			want:     Account{},
			wantErrs: []error{vos.ErrInvalidSecret},
		},
		{
			name:     "should return an error if cpf and secret are invalid",
			cpf:      "12345678900",
			secret:   "invalid",
			want:     Account{},
			wantErrs: []error{vos.ErrInvalidCPF, vos.ErrInvalidSecret},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			acc, err := NewAccount("name", tt.cpf, tt.secret)
			if err != nil {
				for _, e := range tt.wantErrs {
					assert.ErrorIs(t, err, e)
				}
			}

			assert.Equal(t, tt.want.Name, acc.Name)
			assert.Equal(t, tt.want.CPF, acc.CPF)
			assert.Equal(t, tt.want.Secret, acc.Secret)
			assert.Equal(t, tt.want.Balance, acc.Balance)
		})
	}
}
