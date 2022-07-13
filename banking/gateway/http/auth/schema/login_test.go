package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thalissonfelipe/banking/banking/gateway/http/rest"
)

func TestSchema_MapToLoginResponse(t *testing.T) {
	const token = "token"

	got := MapToLoginResponse(token)
	assert.Equal(t, LoginResponse{Token: token}, got)
}

func TestSchema_LoginInput_IsValid(t *testing.T) {
	tests := []struct {
		name    string
		input   LoginInput
		wantErr error
	}{
		{
			name: "should return no error",
			input: LoginInput{
				CPF:    "cpf",
				Secret: "secret",
			},
			wantErr: nil,
		},
		{
			name: "should return error cpf is blank",
			input: LoginInput{
				Secret: "secret",
			},
			wantErr: rest.ErrMissingCPFParameter,
		},
		{
			name: "should return error secret is blank",
			input: LoginInput{
				CPF: "cpf",
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
