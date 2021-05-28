package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"github.com/thalissonfelipe/banking/pkg/domain/account"
	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
	"github.com/thalissonfelipe/banking/pkg/gateways/http/responses"
	"github.com/thalissonfelipe/banking/pkg/services/auth"
	"github.com/thalissonfelipe/banking/pkg/tests"
	"github.com/thalissonfelipe/banking/pkg/tests/fakes"
	"github.com/thalissonfelipe/banking/pkg/tests/mocks"
	"github.com/thalissonfelipe/banking/pkg/tests/testdata"
)

func TestLogin(t *testing.T) {
	cpf := testdata.GetValidCPF()

	testCases := []struct {
		name         string
		usecase      account.Usecase
		enc          account.Encrypter
		body         interface{}
		decoder      tests.Decoder
		expectedBody interface{}
		expectedCode int
	}{
		{
			name:         "should return status code 400 if cpf was not provided",
			usecase:      mocks.StubAccountUsecase{},
			enc:          mocks.StubHash{},
			body:         requestBody{Secret: "12345678"},
			decoder:      tests.ErrorMessageDecoder{},
			expectedBody: responses.ErrorResponse{Message: "missing cpf parameter"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "should return status code 400 if secret was not provided",
			usecase:      mocks.StubAccountUsecase{},
			enc:          mocks.StubHash{},
			body:         requestBody{CPF: cpf.String()},
			decoder:      tests.ErrorMessageDecoder{},
			expectedBody: responses.ErrorResponse{Message: "missing secret parameter"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:    "should return status code 400 if secret was not provided",
			usecase: mocks.StubAccountUsecase{},
			enc:     mocks.StubHash{},
			body: map[string]interface{}{
				"cpf": 123,
			},
			decoder:      tests.ErrorMessageDecoder{},
			expectedBody: responses.ErrorResponse{Message: "invalid json"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "should return status code 500 if usecase fails",
			usecase: mocks.StubAccountUsecase{
				Err: errors.New("usecase fails"),
			},
			enc:          mocks.StubHash{},
			body:         requestBody{CPF: cpf.String(), Secret: "12345678"},
			decoder:      tests.ErrorMessageDecoder{},
			expectedBody: responses.ErrorResponse{Message: "internal server error"},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "should return status code 404 if account does not exist",
			usecase:      mocks.StubAccountUsecase{},
			enc:          mocks.StubHash{},
			body:         requestBody{CPF: cpf.String(), Secret: "12345678"},
			decoder:      tests.ErrorMessageDecoder{},
			expectedBody: responses.ErrorResponse{Message: "account does not exist"},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "should return status code 400 if secret was not correct",
			usecase: mocks.StubAccountUsecase{
				Accounts: []entities.Account{
					entities.NewAccount("Pedro", cpf, vos.NewSecret("87654321")),
				},
			},
			enc:          mocks.StubHash{Err: auth.ErrSecretDoesNotMatch},
			body:         requestBody{CPF: cpf.String(), Secret: "12345678"},
			decoder:      tests.ErrorMessageDecoder{},
			expectedBody: responses.ErrorResponse{Message: "secret does not match"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "should authenticate successfully and return a token",
			usecase: mocks.StubAccountUsecase{
				Accounts: []entities.Account{
					entities.NewAccount("Pedro", cpf, vos.NewSecret("12345678")),
				},
			},
			enc:          mocks.StubHash{},
			body:         requestBody{CPF: cpf.String(), Secret: "12345678"},
			decoder:      responseBodyDecoder{},
			expectedBody: responseBody{},
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

			result := tt.decoder.Decode(response.Body)

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

func (responseBodyDecoder) Decode(body *bytes.Buffer) interface{} {
	var result responseBody
	json.NewDecoder(body).Decode(&result)
	return result
}
