package account

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"github.com/thalissonfelipe/banking/pkg/domain/account/usecase"
	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/gateways/http/account/schemes"
	"github.com/thalissonfelipe/banking/pkg/gateways/http/responses"
	"github.com/thalissonfelipe/banking/pkg/tests"
	"github.com/thalissonfelipe/banking/pkg/tests/fakes"
	"github.com/thalissonfelipe/banking/pkg/tests/mocks"
	"github.com/thalissonfelipe/banking/pkg/tests/testdata"
)

func TestHandler_CreateAccount(t *testing.T) {
	cpf := testdata.GetValidCPF()
	secret := testdata.GetValidSecret()

	testCases := []struct {
		name         string
		repo         *mocks.StubAccountRepository
		enc          *mocks.StubHash
		decoder      tests.Decoder
		body         interface{}
		expectedBody interface{}
		expectedCode int
	}{
		{
			name:         "should return status code 400 if name was not provided",
			repo:         &mocks.StubAccountRepository{},
			enc:          &mocks.StubHash{},
			decoder:      tests.ErrorMessageDecoder{},
			body:         schemes.CreateAccountInput{CPF: cpf.String(), Secret: secret.String()},
			expectedBody: responses.ErrorResponse{Message: "missing name parameter"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "should return status code 400 if cpf was not provided",
			repo:         &mocks.StubAccountRepository{},
			enc:          &mocks.StubHash{},
			decoder:      tests.ErrorMessageDecoder{},
			body:         schemes.CreateAccountInput{Name: "Pedro", Secret: secret.String()},
			expectedBody: responses.ErrorResponse{Message: "missing cpf parameter"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "should return status code 400 if secret was not provided",
			repo:         &mocks.StubAccountRepository{},
			enc:          &mocks.StubHash{},
			decoder:      tests.ErrorMessageDecoder{},
			body:         schemes.CreateAccountInput{Name: "Pedro", CPF: cpf.String()},
			expectedBody: responses.ErrorResponse{Message: "missing secret parameter"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "should return status code 400 if an invalid json was provided",
			repo:         &mocks.StubAccountRepository{},
			enc:          &mocks.StubHash{},
			decoder:      tests.ErrorMessageDecoder{},
			body:         map[string]interface{}{"name": 123456},
			expectedBody: responses.ErrorResponse{Message: "invalid json"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "should return status code 400 if cpf provided is not valid",
			repo:         &mocks.StubAccountRepository{},
			enc:          &mocks.StubHash{},
			decoder:      tests.ErrorMessageDecoder{},
			body:         schemes.CreateAccountInput{Name: "Pedro", CPF: "123.456.789-00", Secret: secret.String()},
			expectedBody: responses.ErrorResponse{Message: "invalid cpf"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "should return status code 400 if secret provided is not valid",
			repo:         &mocks.StubAccountRepository{},
			enc:          &mocks.StubHash{},
			decoder:      tests.ErrorMessageDecoder{},
			body:         schemes.CreateAccountInput{Name: "Pedro", CPF: cpf.String(), Secret: "12345678"},
			expectedBody: responses.ErrorResponse{Message: "invalid secret"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "should return status code 409 if cpf already exists",
			repo: &mocks.StubAccountRepository{
				Accounts: []entities.Account{entities.NewAccount(
					"Junior", cpf, secret,
				)},
			},
			enc:          &mocks.StubHash{},
			decoder:      tests.ErrorMessageDecoder{},
			body:         schemes.CreateAccountInput{Name: "Pedro", CPF: cpf.String(), Secret: secret.String()},
			expectedBody: responses.ErrorResponse{Message: "account already exists"},
			expectedCode: http.StatusConflict,
		},
		{
			name: "should return status code 500 if usecase fails to create account",
			repo: &mocks.StubAccountRepository{
				Err: errors.New("usecase error"),
			},
			enc:          &mocks.StubHash{},
			decoder:      tests.ErrorMessageDecoder{},
			body:         schemes.CreateAccountInput{Name: "Pedro", CPF: cpf.String(), Secret: secret.String()},
			expectedBody: responses.ErrorResponse{Message: "internal server error"},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "should return status code 500 if hash fails to hash secret",
			repo: &mocks.StubAccountRepository{
				Err: errors.New("usecase error"),
			},
			enc:          &mocks.StubHash{Err: errors.New("hash error")},
			decoder:      tests.ErrorMessageDecoder{},
			body:         schemes.CreateAccountInput{Name: "Pedro", CPF: cpf.String(), Secret: secret.String()},
			expectedBody: responses.ErrorResponse{Message: "internal server error"},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "should return status code 201 and created account",
			repo:         &mocks.StubAccountRepository{},
			enc:          &mocks.StubHash{},
			decoder:      createdAccountDecoder{},
			body:         schemes.CreateAccountInput{Name: "Pedro", CPF: cpf.String(), Secret: secret.String()},
			expectedBody: schemes.CreateAccountResponse{Name: "Pedro", CPF: cpf.String(), Balance: 0},
			expectedCode: http.StatusCreated,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			usecase := usecase.NewAccountUsecase(tt.repo, tt.enc)
			handler := NewHandler(r, usecase)

			request := fakes.FakeRequest(http.MethodPost, "/accounts", tt.body)
			response := httptest.NewRecorder()

			http.HandlerFunc(handler.CreateAccount).ServeHTTP(response, request)

			result := tt.decoder.Decode(response.Body)

			assert.Equal(t, tt.expectedBody, result)
			assert.Equal(t, tt.expectedCode, response.Code)
		})
	}
}

type createdAccountDecoder struct{}

func (createdAccountDecoder) Decode(body *bytes.Buffer) interface{} {
	var result schemes.CreateAccountResponse
	json.NewDecoder(body).Decode(&result)
	return result
}
