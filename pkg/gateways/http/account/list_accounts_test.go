package account

import (
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
	"github.com/thalissonfelipe/banking/pkg/tests/fakes"
	"github.com/thalissonfelipe/banking/pkg/tests/mocks"
)

func TestListAccounts(t *testing.T) {
	acc := entities.NewAccount("Pedro", "123.456.789-00", "12345678")
	testCases := []struct {
		name         string
		repoSetup    *mocks.StubAccountRepository
		expectedBody interface{}
		expectedCode int
	}{
		{
			name:         "should return 200 and an empty slice of accounts",
			repoSetup:    &mocks.StubAccountRepository{},
			expectedBody: []AccountResponse{},
			expectedCode: http.StatusOK,
		},
		{
			name: "should return 200 and an slice of accounts",
			repoSetup: &mocks.StubAccountRepository{
				Accounts: []entities.Account{acc},
			},
			expectedBody: []AccountResponse{convertAccountToAccountResponse(acc)},
			expectedCode: http.StatusOK,
		},
		{
			name:         "should return 500 and error message if something went wrong",
			repoSetup:    &mocks.StubAccountRepository{Err: errors.New("failed to list accounts")},
			expectedBody: responses.ErrorResponse{Message: "Internal Error."},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			accUseCase := usecase.NewAccountUseCase(tt.repoSetup, nil)
			r := mux.NewRouter()
			handler := NewHandler(r, accUseCase)

			request := fakes.FakeRequest(http.MethodGet, "/accounts", nil)
			response := httptest.NewRecorder()

			http.HandlerFunc(handler.ListAccounts).ServeHTTP(response, request)

			assert.Equal(t, tt.expectedCode, response.Code)
			assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
			assertResponseBody(t, tt.expectedBody, response)
		})
	}
}

func assertResponseBody(t *testing.T, expected interface{}, response *httptest.ResponseRecorder) {
	t.Helper()

	expectedBytes, err := json.Marshal(expected)
	if err != nil {
		t.Errorf("could not marshall expected interface")
	}

	var result interface{}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		t.Errorf("could not decode response body")
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		t.Errorf("could not marshall response body")
	}

	assert.ObjectsAreEqualValues(expectedBytes, resultBytes)
}
