package account

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/banking/domain/entities"
	"github.com/thalissonfelipe/banking/banking/gateway/http/account/schemes"
	"github.com/thalissonfelipe/banking/banking/gateway/http/rest"
	"github.com/thalissonfelipe/banking/banking/tests"
	"github.com/thalissonfelipe/banking/banking/tests/fakes"
	"github.com/thalissonfelipe/banking/banking/tests/mocks"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
)

func TestHandler_CreateAccount(t *testing.T) {
	cpf := testdata.GetValidCPF()
	secret := testdata.GetValidSecret()

	testCases := []struct {
		name         string
		usecase      *mocks.AccountUsecaseMock
		enc          *mocks.HashMock
		decoder      tests.Decoder
		body         interface{}
		expectedBody interface{}
		expectedCode int
	}{
		{
			name:         "should return status code 400 if name was not provided",
			usecase:      &mocks.AccountUsecaseMock{},
			enc:          &mocks.HashMock{},
			decoder:      tests.ErrorMessageDecoder{},
			body:         schemes.CreateAccountInput{CPF: cpf.String(), Secret: secret.String()},
			expectedBody: rest.ErrorResponse{Message: "missing name parameter"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "should return status code 400 if cpf was not provided",
			usecase:      &mocks.AccountUsecaseMock{},
			enc:          &mocks.HashMock{},
			decoder:      tests.ErrorMessageDecoder{},
			body:         schemes.CreateAccountInput{Name: "Pedro", Secret: secret.String()},
			expectedBody: rest.ErrorResponse{Message: "missing cpf parameter"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "should return status code 400 if secret was not provided",
			usecase:      &mocks.AccountUsecaseMock{},
			enc:          &mocks.HashMock{},
			decoder:      tests.ErrorMessageDecoder{},
			body:         schemes.CreateAccountInput{Name: "Pedro", CPF: cpf.String()},
			expectedBody: rest.ErrorResponse{Message: "missing secret parameter"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "should return status code 400 if an invalid json was provided",
			usecase:      &mocks.AccountUsecaseMock{},
			enc:          &mocks.HashMock{},
			decoder:      tests.ErrorMessageDecoder{},
			body:         map[string]interface{}{"name": 123456},
			expectedBody: rest.ErrorResponse{Message: "invalid json"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "should return status code 400 if cpf provided is not valid",
			usecase:      &mocks.AccountUsecaseMock{},
			enc:          &mocks.HashMock{},
			decoder:      tests.ErrorMessageDecoder{},
			body:         schemes.CreateAccountInput{Name: "Pedro", CPF: "123.456.789-00", Secret: secret.String()},
			expectedBody: rest.ErrorResponse{Message: "invalid cpf"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "should return status code 400 if secret provided is not valid",
			usecase:      &mocks.AccountUsecaseMock{},
			enc:          &mocks.HashMock{},
			decoder:      tests.ErrorMessageDecoder{},
			body:         schemes.CreateAccountInput{Name: "Pedro", CPF: cpf.String(), Secret: "12345678"},
			expectedBody: rest.ErrorResponse{Message: "invalid secret"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "should return status code 409 if cpf already exists",
			usecase: &mocks.AccountUsecaseMock{
				Accounts: []entities.Account{entities.NewAccount(
					"Junior", cpf, secret,
				)},
			},
			enc:          &mocks.HashMock{},
			decoder:      tests.ErrorMessageDecoder{},
			body:         schemes.CreateAccountInput{Name: "Pedro", CPF: cpf.String(), Secret: secret.String()},
			expectedBody: rest.ErrorResponse{Message: "account already exists"},
			expectedCode: http.StatusConflict,
		},
		{
			name:         "should return status code 500 if usecase fails to create account",
			usecase:      &mocks.AccountUsecaseMock{Err: testdata.ErrUsecaseFails},
			enc:          &mocks.HashMock{},
			decoder:      tests.ErrorMessageDecoder{},
			body:         schemes.CreateAccountInput{Name: "Pedro", CPF: cpf.String(), Secret: secret.String()},
			expectedBody: rest.ErrorResponse{Message: "internal server error"},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "should return status code 201 and created account",
			usecase:      &mocks.AccountUsecaseMock{},
			enc:          &mocks.HashMock{},
			decoder:      createdAccountDecoder{},
			body:         schemes.CreateAccountInput{Name: "Pedro", CPF: cpf.String(), Secret: secret.String()},
			expectedBody: schemes.CreateAccountResponse{Name: "Pedro", CPF: cpf.String(), Balance: 0},
			expectedCode: http.StatusCreated,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			handler := NewHandler(r, tt.usecase)

			request := fakes.FakeRequest(http.MethodPost, "/accounts", tt.body)
			response := httptest.NewRecorder()

			http.HandlerFunc(handler.CreateAccount).ServeHTTP(response, request)

			result := tt.decoder.Decode(t, response.Body)

			assert.Equal(t, tt.expectedBody, result)
			assert.Equal(t, tt.expectedCode, response.Code)
		})
	}
}

type createdAccountDecoder struct{}

func (createdAccountDecoder) Decode(t *testing.T, body *bytes.Buffer) interface{} {
	var result schemes.CreateAccountResponse

	err := json.NewDecoder(body).Decode(&result)
	require.NoError(t, err)

	return result
}
