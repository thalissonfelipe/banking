package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thalissonfelipe/banking/banking/domain/vos"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
)

func TestNewAccount(t *testing.T) {
	cpf := testdata.GetValidCPF()
	secret := testdata.GetValidSecret()

	tests := []struct {
		name    string
		cpf     string
		secret  string
		want    Account
		wantErr error
	}{
		{
			name:   "should create an account successfully",
			cpf:    cpf.String(),
			secret: secret.String(),
			want: Account{
				Name:    "name",
				CPF:     cpf,
				Secret:  secret,
				Balance: 0,
			},
			wantErr: nil,
		},
		{
			name:    "should return an error if cpf is invalid",
			cpf:     "12345678900",
			secret:  secret.String(),
			want:    Account{},
			wantErr: vos.ErrInvalidCPF,
		},
		{
			name:    "should return an error if secret is invalid",
			cpf:     cpf.String(),
			secret:  "invalid",
			want:    Account{},
			wantErr: vos.ErrInvalidSecret,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acc, err := NewAccount("name", tt.cpf, tt.secret)
			assert.ErrorIs(t, err, tt.wantErr)

			assert.Equal(t, tt.want.Name, acc.Name)
			assert.Equal(t, tt.want.CPF, acc.CPF)
			assert.Equal(t, tt.want.Secret, acc.Secret)
			assert.Equal(t, tt.want.Balance, acc.Balance)
		})
	}
}
