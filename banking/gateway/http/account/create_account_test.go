package account

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/domain/usecases"
	"github.com/thalissonfelipe/banking/banking/gateway/http/account/schema"
	"github.com/thalissonfelipe/banking/banking/gateway/http/rest"
	"github.com/thalissonfelipe/banking/banking/tests/fakes"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
)

func TestAccountHandler_CreateAccount(t *testing.T) {
	t.Parallel()

	cpf := testdata.CPF()
	secret := testdata.Secret()

	acc, err := entity.NewAccount("name", cpf.String(), secret.String())
	require.NoError(t, err)

	acc.CreatedAt = time.Now()

	testCases := []struct {
		name     string
		usecase  usecases.Account
		body     interface{}
		wantBody interface{}
		wantCode int
	}{
		{
			name: "should return status code 201 if acc was created successfully",
			usecase: &UsecaseMock{
				CreateAccountFunc: func(_ context.Context, account *entity.Account) error {
					*account = acc

					return nil
				},
			},
			body: schema.CreateAccountInput{Name: "name", CPF: cpf.String(), Secret: secret.String()},
			wantBody: schema.CreateAccountResponse{
				ID:        acc.ID.String(),
				Name:      "name",
				CPF:       cpf.String(),
				Balance:   100,
				CreatedAt: acc.CreatedAt.UTC().Format(time.RFC3339),
			},
			wantCode: http.StatusCreated,
		},
		{
			name:    "should return status code 400 if name was not provided",
			usecase: &UsecaseMock{},
			body:    schema.CreateAccountInput{CPF: cpf.String(), Secret: secret.String()},
			wantBody: rest.Error{
				Error: "invalid request body",
				Details: []rest.ErrorDetail{
					{
						Location: "body.name",
						Message:  "missing parameter",
					},
				},
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "should return status code 400 if an invalid json was provided",
			usecase:  &UsecaseMock{},
			body:     map[string]interface{}{"name": 123456},
			wantBody: rest.Error{Error: "invalid request body"},
			wantCode: http.StatusBadRequest,
		},
		{
			name:    "should return status code 400 if cpf provided is not valid",
			usecase: &UsecaseMock{},
			body:    schema.CreateAccountInput{Name: "name", CPF: "123.456.789-00", Secret: secret.String()},
			wantBody: rest.Error{
				Error: "invalid request body",
				Details: []rest.ErrorDetail{
					{
						Location: "body.cpf",
						Message:  "invalid cpf",
					},
				},
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "should return status code 409 if account already exists",
			usecase: &UsecaseMock{
				CreateAccountFunc: func(ctx context.Context, account *entity.Account) error {
					return entity.ErrAccountAlreadyExists
				},
			},
			body:     schema.CreateAccountInput{Name: "name", CPF: cpf.String(), Secret: secret.String()},
			wantBody: rest.Error{Error: "account already exists"},
			wantCode: http.StatusConflict,
		},
		{
			name: "should return status code 500 if usecase fails to create an account",
			usecase: &UsecaseMock{
				CreateAccountFunc: func(ctx context.Context, account *entity.Account) error {
					return assert.AnError
				},
			},
			body:     schema.CreateAccountInput{Name: "name", CPF: cpf.String(), Secret: secret.String()},
			wantBody: rest.Error{Error: "internal server error"},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range testCases {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := NewHandler(tt.usecase)

			request := fakes.FakeRequest(http.MethodPost, "/accounts", tt.body)
			response := httptest.NewRecorder()

			rest.Wrap(handler.CreateAccount).ServeHTTP(response, request)

			want, err := json.Marshal(tt.wantBody)
			require.NoError(t, err)

			assert.Equal(t, tt.wantCode, response.Code)
			assert.JSONEq(t, string(want), response.Body.String())
			assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
		})
	}
}
