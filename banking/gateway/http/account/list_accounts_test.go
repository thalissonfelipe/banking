package account

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

	"github.com/thalissonfelipe/banking/banking/domain/account"
	"github.com/thalissonfelipe/banking/banking/domain/entities"
	"github.com/thalissonfelipe/banking/banking/gateway/http/account/schemes"
	"github.com/thalissonfelipe/banking/banking/gateway/http/rest"
	"github.com/thalissonfelipe/banking/banking/tests"
	"github.com/thalissonfelipe/banking/banking/tests/fakes"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
)

func TestAccountHandler_ListAccounts(t *testing.T) {
	acc, err := entities.NewAccount("name", testdata.GetValidCPF().String(), testdata.GetValidSecret().String())
	require.NoError(t, err)

	testCases := []struct {
		name         string
		usecase      account.Usecase
		expectedBody interface{}
		decoder      tests.Decoder
		expectedCode int
	}{
		{
			name: "should return a list of accounts successfully",
			usecase: &UsecaseMock{
				ListAccountsFunc: func(context.Context) ([]entities.Account, error) {
					return []entities.Account{acc}, nil
				},
			},
			expectedBody: []schemes.AccountListResponse{convertAccountToAccountListResponse(acc)},
			expectedCode: http.StatusOK,
			decoder:      listAccountsSuccessDecoder{},
		},
		{
			name: "should return 500 if usecase fails",
			usecase: &UsecaseMock{
				ListAccountsFunc: func(context.Context) ([]entities.Account, error) {
					return nil, assert.AnError
				},
			},
			expectedBody: rest.ErrorResponse{Message: "internal server error"},
			decoder:      tests.ErrorMessageDecoder{},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			handler := NewHandler(r, tt.usecase)

			request := fakes.FakeRequest(http.MethodGet, "/accounts", nil)
			response := httptest.NewRecorder()

			http.HandlerFunc(handler.ListAccounts).ServeHTTP(response, request)

			result := tt.decoder.Decode(t, response.Body)

			assert.Equal(t, tt.expectedCode, response.Code)
			assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
			assert.Equal(t, tt.expectedBody, result)
		})
	}
}

type listAccountsSuccessDecoder struct{}

func (listAccountsSuccessDecoder) Decode(t *testing.T, body *bytes.Buffer) interface{} {
	var result []schemes.AccountListResponse

	err := json.NewDecoder(body).Decode(&result)
	require.NoError(t, err)

	return result
}
