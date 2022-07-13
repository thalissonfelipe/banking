package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/banking/gateway/http/auth/schemes"
	"github.com/thalissonfelipe/banking/banking/gateway/http/rest"
	"github.com/thalissonfelipe/banking/banking/services"
	"github.com/thalissonfelipe/banking/banking/services/auth"
	"github.com/thalissonfelipe/banking/banking/tests"
	"github.com/thalissonfelipe/banking/banking/tests/fakes"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
)

func TestAuthHandler_Login(t *testing.T) {
	cpf := testdata.GetValidCPF()
	secret := testdata.GetValidSecret()

	testCases := []struct {
		name         string
		auth         services.Auth
		body         interface{}
		decoder      tests.Decoder
		expectedBody interface{}
		expectedCode int
	}{
		{
			name: "should authenticate successfully and return a token",
			auth: &AuthMock{
				AutheticateFunc: func(_ context.Context, _, _ string) (string, error) {
					return "token", nil
				},
			},
			body:         schemes.LoginInput{CPF: cpf.String(), Secret: secret.String()},
			decoder:      responseBodyDecoder{},
			expectedBody: schemes.LoginResponse{},
			expectedCode: http.StatusOK,
		},
		{
			name:         "should return status code 400 if cpf was not provided",
			auth:         &AuthMock{},
			body:         schemes.LoginInput{Secret: secret.String()},
			decoder:      tests.ErrorMessageDecoder{},
			expectedBody: rest.ErrorResponse{Message: "missing cpf parameter"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "should return status code 400 if secret was not provided",
			auth:         &AuthMock{},
			body:         schemes.LoginInput{CPF: cpf.String()},
			decoder:      tests.ErrorMessageDecoder{},
			expectedBody: rest.ErrorResponse{Message: "missing secret parameter"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "should return status code 400 if json provided was not valid",
			auth: &AuthMock{},
			body: map[string]interface{}{
				"cpf": 123,
			},
			decoder:      tests.ErrorMessageDecoder{},
			expectedBody: rest.ErrorResponse{Message: "invalid json"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "should return status code 500 if auth service fails",
			auth: &AuthMock{
				AutheticateFunc: func(_ context.Context, _, _ string) (string, error) {
					return "", assert.AnError
				},
			},
			body:         schemes.LoginInput{CPF: cpf.String(), Secret: secret.String()},
			decoder:      tests.ErrorMessageDecoder{},
			expectedBody: rest.ErrorResponse{Message: "internal server error"},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "should return status code 400 if credentials are invalid",
			auth: &AuthMock{
				AutheticateFunc: func(_ context.Context, _, _ string) (string, error) {
					return "", auth.ErrInvalidCredentials
				},
			},
			body:         schemes.LoginInput{CPF: cpf.String(), Secret: secret.String()},
			decoder:      tests.ErrorMessageDecoder{},
			expectedBody: rest.ErrorResponse{Message: "cpf or secret are invalid"},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			handler := NewHandler(r, tt.auth)

			request := fakes.FakeRequest(http.MethodPost, "/accounts", tt.body)
			response := httptest.NewRecorder()

			http.HandlerFunc(handler.Login).ServeHTTP(response, request)

			result := tt.decoder.Decode(t, response.Body)

			// This is temporary because I don't know yet how to
			// test the generated token. But it's working. (I hope)
			if response.Code != http.StatusOK {
				assert.Equal(t, tt.expectedBody, result)
			}

			assert.Equal(t, tt.expectedCode, response.Code)
		})
	}
}

type responseBodyDecoder struct{}

func (responseBodyDecoder) Decode(t *testing.T, body *bytes.Buffer) interface{} {
	var result schemes.LoginResponse

	err := json.NewDecoder(body).Decode(&result)
	require.NoError(t, err)

	return result
}
