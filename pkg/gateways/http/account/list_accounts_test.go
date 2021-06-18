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

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/gateways/http/account/schemes"
	"github.com/thalissonfelipe/banking/pkg/gateways/http/rest"
	"github.com/thalissonfelipe/banking/pkg/tests"
	"github.com/thalissonfelipe/banking/pkg/tests/fakes"
	"github.com/thalissonfelipe/banking/pkg/tests/mocks"
	"github.com/thalissonfelipe/banking/pkg/tests/testdata"
)

func TestHandler_ListAccounts(t *testing.T) {
	acc := entities.NewAccount("Pedro", testdata.GetValidCPF(), testdata.GetValidSecret())

	testCases := []struct {
		name         string
		usecase      *mocks.AccountUsecaseMock
		expectedBody interface{}
		decoder      tests.Decoder
		expectedCode int
	}{
		{
			name:         "should return 200 and an empty slice of accounts",
			usecase:      &mocks.AccountUsecaseMock{},
			expectedBody: []schemes.AccountListResponse{},
			expectedCode: http.StatusOK,
			decoder:      listAccountsSuccessDecoder{},
		},
		{
			name: "should return 200 and an slice of accounts",
			usecase: &mocks.AccountUsecaseMock{
				Accounts: []entities.Account{acc},
			},
			expectedBody: []schemes.AccountListResponse{convertAccountToAccountListResponse(acc)},
			decoder:      listAccountsSuccessDecoder{},
			expectedCode: http.StatusOK,
		},
		{
			name:         "should return 500 and error message if something went wrong",
			usecase:      &mocks.AccountUsecaseMock{Err: testdata.ErrRepositoryFailsToFetch},
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
