package account

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/domain/usecases"
	"github.com/thalissonfelipe/banking/banking/gateway/http/account/schema"
	"github.com/thalissonfelipe/banking/banking/gateway/http/rest"
	"github.com/thalissonfelipe/banking/banking/tests"
	"github.com/thalissonfelipe/banking/banking/tests/fakes"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
)

func TestAccountHandler_CreateAccount(t *testing.T) {
	cpf := testdata.GetValidCPF()
	secret := testdata.GetValidSecret()

	acc, err := entity.NewAccount("name", cpf.String(), secret.String())
	require.NoError(t, err)

	acc.CreatedAt = time.Now()

	testCases := []struct {
		name         string
		usecase      usecases.Account
		decoder      tests.Decoder
		body         interface{}
		expectedBody interface{}
		expectedCode int
	}{
		{
			name: "should return status code 201 if acc was created successfully",
			usecase: &UsecaseMock{
				CreateAccountFunc: func(_ context.Context, account *entity.Account) error {
					*account = acc
					return nil
				},
			},
			decoder: createdAccountDecoder{},
			body:    schema.CreateAccountInput{Name: "name", CPF: cpf.String(), Secret: secret.String()},
			expectedBody: schema.CreateAccountResponse{
				Name:      "name",
				CPF:       cpf.String(),
				Balance:   100,
				CreatedAt: acc.CreatedAt.UTC().Format(time.RFC3339),
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:         "should return status code 400 if name was not provided",
			usecase:      &UsecaseMock{},
			decoder:      tests.ErrorMessageDecoder{},
			body:         schema.CreateAccountInput{CPF: cpf.String(), Secret: secret.String()},
			expectedBody: rest.ErrorResponse{Message: "missing name parameter"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "should return status code 400 if cpf was not provided",
			usecase:      &UsecaseMock{},
			decoder:      tests.ErrorMessageDecoder{},
			body:         schema.CreateAccountInput{Name: "name", Secret: secret.String()},
			expectedBody: rest.ErrorResponse{Message: "missing cpf parameter"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "should return status code 400 if secret was not provided",
			usecase:      &UsecaseMock{},
			decoder:      tests.ErrorMessageDecoder{},
			body:         schema.CreateAccountInput{Name: "name", CPF: cpf.String()},
			expectedBody: rest.ErrorResponse{Message: "missing secret parameter"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "should return status code 400 if an invalid json was provided",
			usecase:      &UsecaseMock{},
			decoder:      tests.ErrorMessageDecoder{},
			body:         map[string]interface{}{"name": 123456},
			expectedBody: rest.ErrorResponse{Message: "invalid json"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "should return status code 400 if cpf provided is not valid",
			usecase:      &UsecaseMock{},
			decoder:      tests.ErrorMessageDecoder{},
			body:         schema.CreateAccountInput{Name: "name", CPF: "123.456.789-00", Secret: secret.String()},
			expectedBody: rest.ErrorResponse{Message: "invalid cpf"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "should return status code 400 if secret provided is not valid",
			usecase:      &UsecaseMock{},
			decoder:      tests.ErrorMessageDecoder{},
			body:         schema.CreateAccountInput{Name: "name", CPF: cpf.String(), Secret: "12345678"},
			expectedBody: rest.ErrorResponse{Message: "invalid secret"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "should return status code 409 if account already exists",
			usecase: &UsecaseMock{
				CreateAccountFunc: func(ctx context.Context, account *entity.Account) error {
					return entity.ErrAccountAlreadyExists
				},
			},
			decoder:      tests.ErrorMessageDecoder{},
			body:         schema.CreateAccountInput{Name: "name", CPF: cpf.String(), Secret: secret.String()},
			expectedBody: rest.ErrorResponse{Message: "account already exists"},
			expectedCode: http.StatusConflict,
		},
		{
			name: "should return status code 500 if usecase fails to create an account",
			usecase: &UsecaseMock{
				CreateAccountFunc: func(ctx context.Context, account *entity.Account) error {
					return assert.AnError
				},
			},
			decoder:      tests.ErrorMessageDecoder{},
			body:         schema.CreateAccountInput{Name: "name", CPF: cpf.String(), Secret: secret.String()},
			expectedBody: rest.ErrorResponse{Message: "internal server error"},
			expectedCode: http.StatusInternalServerError,
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
	var result schema.CreateAccountResponse

	err := json.NewDecoder(body).Decode(&result)
	require.NoError(t, err)

	return result
}
