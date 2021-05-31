package account

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/thalissonfelipe/banking/pkg/domain/account/usecase"
	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/gateways/http/responses"
	"github.com/thalissonfelipe/banking/pkg/tests"
	"github.com/thalissonfelipe/banking/pkg/tests/fakes"
	"github.com/thalissonfelipe/banking/pkg/tests/mocks"
	"github.com/thalissonfelipe/banking/pkg/tests/testdata"
)

func TestHandler_ListAccounts(t *testing.T) {
	acc := entities.NewAccount("Pedro", testdata.GetValidCPF(), testdata.GetValidSecret())

	testCases := []struct {
		name         string
		repoSetup    *mocks.StubAccountRepository
		expectedBody interface{}
		decoder      tests.Decoder
		expectedCode int
	}{
		{
			name:         "should return 200 and an empty slice of accounts",
			repoSetup:    &mocks.StubAccountRepository{},
			expectedBody: []accountResponse{},
			expectedCode: http.StatusOK,
			decoder:      listAccountsSuccessDecoder{},
		},
		{
			name: "should return 200 and an slice of accounts",
			repoSetup: &mocks.StubAccountRepository{
				Accounts: []entities.Account{acc},
			},
			expectedBody: []accountResponse{convertAccountToAccountResponse(acc)},
			decoder:      listAccountsSuccessDecoder{},
			expectedCode: http.StatusOK,
		},
		{
			name:         "should return 500 and error message if something went wrong",
			repoSetup:    &mocks.StubAccountRepository{Err: errors.New("failed to list accounts")},
			expectedBody: responses.ErrorResponse{Message: "internal server error"},
			decoder:      tests.ErrorMessageDecoder{},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			accUseCase := usecase.NewAccountUsecase(tt.repoSetup, nil)
			r := mux.NewRouter()
			handler := NewHandler(r, accUseCase)

			request := fakes.FakeRequest(http.MethodGet, "/accounts", nil)
			response := httptest.NewRecorder()

			http.HandlerFunc(handler.ListAccounts).ServeHTTP(response, request)

			result := tt.decoder.Decode(response.Body)

			assert.Equal(t, tt.expectedCode, response.Code)
			assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
			assert.Equal(t, tt.expectedBody, result)
		})
	}
}

type listAccountsSuccessDecoder struct{}

func (listAccountsSuccessDecoder) Decode(body *bytes.Buffer) interface{} {
	var result []accountResponse
	json.NewDecoder(body).Decode(&result)
	return result
}
