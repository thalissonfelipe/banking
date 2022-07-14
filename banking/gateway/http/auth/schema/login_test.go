package schema

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		wantErr rest.ValidationErrors
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
			name: "should return error if cpf is blank",
			input: LoginInput{
				Secret: "secret",
			},
			wantErr: rest.ValidationErrors{ErrMissingCPFParameter},
		},
		{
			name: "should return error if secret is blank",
			input: LoginInput{
				CPF: "cpf",
			},
			wantErr: rest.ValidationErrors{ErrMissingSecretParameter},
		},
		{
			name:    "should return error if cpf and secret are blank",
			input:   LoginInput{},
			wantErr: rest.ValidationErrors{ErrMissingCPFParameter, ErrMissingSecretParameter},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
