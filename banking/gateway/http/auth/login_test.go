package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/banking/domain/usecases"
	"github.com/thalissonfelipe/banking/banking/gateway/http/auth/schema"
	"github.com/thalissonfelipe/banking/banking/gateway/http/rest"
	"github.com/thalissonfelipe/banking/banking/tests/fakes"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
)

func TestAuthHandler_Login(t *testing.T) {
	t.Parallel()

	cpf := testdata.CPF()
	secret := testdata.Secret()

	testCases := []struct {
		name     string
		auth     usecases.Auth
		body     interface{}
		wantBody interface{}
		wantCode int
	}{
		{
			name: "should authenticate successfully and return a token",
			auth: &UsecaseMock{
				AutheticateFunc: func(_ context.Context, _, _ string) (string, error) {
					return "token", nil
				},
			},
			body:     schema.LoginInput{CPF: cpf.String(), Secret: secret.String()},
			wantBody: schema.LoginResponse{Token: "token"},
			wantCode: http.StatusOK,
		},
		{
			name: "should return status code 400 if cpf was not provided",
			auth: &UsecaseMock{},
			body: schema.LoginInput{Secret: secret.String()},
			wantBody: rest.Error{
				Error: "invalid request body",
				Details: []rest.ErrorDetail{
					{
						Location: "body.cpf",
						Message:  "missing parameter",
					},
				},
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "should return status code 400 if json provided was not valid",
			auth: &UsecaseMock{},
			body: map[string]interface{}{
				"cpf": 123,
			},
			wantBody: rest.Error{Error: "invalid request body"},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "should return status code 400 if credentials are invalid",
			auth: &UsecaseMock{
				AutheticateFunc: func(_ context.Context, _, _ string) (string, error) {
					return "", usecases.ErrInvalidCredentials
				},
			},
			body:     schema.LoginInput{CPF: cpf.String(), Secret: secret.String()},
			wantBody: rest.Error{Error: "invalid credentials"},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "should return status code 500 if auth service fails",
			auth: &UsecaseMock{
				AutheticateFunc: func(_ context.Context, _, _ string) (string, error) {
					return "", assert.AnError
				},
			},
			body:     schema.LoginInput{CPF: cpf.String(), Secret: secret.String()},
			wantBody: rest.Error{Error: "internal server error"},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range testCases {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := NewHandler(tt.auth)

			request := fakes.FakeRequest(http.MethodPost, "/login", tt.body)
			response := httptest.NewRecorder()

			rest.Wrap(handler.Login).ServeHTTP(response, request)

			want, err := json.Marshal(tt.wantBody)
			require.NoError(t, err)

			assert.Equal(t, tt.wantCode, response.Code)
			assert.JSONEq(t, string(want), response.Body.String())
			assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
		})
	}
}
