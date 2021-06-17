package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/pkg/domain/account"
	"github.com/thalissonfelipe/banking/pkg/domain/encrypter"
	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/gateways/http/auth/schemes"
	"github.com/thalissonfelipe/banking/pkg/gateways/http/responses"
	"github.com/thalissonfelipe/banking/pkg/services/auth"
	"github.com/thalissonfelipe/banking/pkg/tests"
	"github.com/thalissonfelipe/banking/pkg/tests/fakes"
	"github.com/thalissonfelipe/banking/pkg/tests/mocks"
	"github.com/thalissonfelipe/banking/pkg/tests/testdata"
)

func TestLogin(t *testing.T) {
	cpf := testdata.GetValidCPF()
	secret := testdata.GetValidSecret()

	testCases := []struct {
		name         string
		usecase      account.Usecase
		enc          encrypter.Encrypter
		body         interface{}
		decoder      tests.Decoder
		expectedBody interface{}
		expectedCode int
	}{
		{
			name:         "should return status code 400 if cpf was not provided",
			usecase:      &mocks.StubAccountUsecase{},
			enc:          mocks.StubHash{},
			body:         schemes.LoginInput{Secret: secret.String()},
			decoder:      tests.ErrorMessageDecoder{},
			expectedBody: responses.ErrorResponse{Message: "missing cpf parameter"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "should return status code 400 if secret was not provided",
			usecase:      &mocks.StubAccountUsecase{},
			enc:          mocks.StubHash{},
			body:         schemes.LoginInput{CPF: cpf.String()},
			decoder:      tests.ErrorMessageDecoder{},
			expectedBody: responses.ErrorResponse{Message: "missing secret parameter"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:    "should return status code 400 if json provided was not valid",
			usecase: &mocks.StubAccountUsecase{},
			enc:     mocks.StubHash{},
			body: map[string]interface{}{
				"cpf": 123,
			},
			decoder:      tests.ErrorMessageDecoder{},
			expectedBody: responses.ErrorResponse{Message: "invalid json"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "should return status code 500 if usecase fails",
			usecase:      &mocks.StubAccountUsecase{Err: testdata.ErrUsecaseFails},
			enc:          mocks.StubHash{},
			body:         schemes.LoginInput{CPF: cpf.String(), Secret: secret.String()},
			decoder:      tests.ErrorMessageDecoder{},
			expectedBody: responses.ErrorResponse{Message: "internal server error"},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "should return status code 400 if account does not exist",
			usecase:      &mocks.StubAccountUsecase{},
			enc:          mocks.StubHash{},
			body:         schemes.LoginInput{CPF: cpf.String(), Secret: secret.String()},
			decoder:      tests.ErrorMessageDecoder{},
			expectedBody: responses.ErrorResponse{Message: "cpf or secret are invalid"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "should return status code 400 if secret was not correct",
			usecase: &mocks.StubAccountUsecase{
				Accounts: []entities.Account{
					entities.NewAccount("Pedro", cpf, secret),
				},
			},
			enc:          mocks.StubHash{Err: auth.ErrInvalidCredentials},
			body:         schemes.LoginInput{CPF: cpf.String(), Secret: "12345678"},
			decoder:      tests.ErrorMessageDecoder{},
			expectedBody: responses.ErrorResponse{Message: "cpf or secret are invalid"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "should authenticate successfully and return a token",
			usecase: &mocks.StubAccountUsecase{
				Accounts: []entities.Account{
					entities.NewAccount("Pedro", cpf, secret),
				},
			},
			enc:          mocks.StubHash{},
			body:         schemes.LoginInput{CPF: cpf.String(), Secret: secret.String()},
			decoder:      responseBodyDecoder{},
			expectedBody: schemes.LoginResponse{},
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			authService := auth.NewAuth(tt.usecase, tt.enc)
			handler := NewHandler(r, authService)

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
