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
		repoSetup    func() mocks.StubAccountRepository
		expectedBody func() interface{}
		expectedCode int
	}{
		{
			name: "should return 200 and an empty slice of accounts",
			repoSetup: func() mocks.StubAccountRepository {
				return mocks.StubAccountRepository{}
			},
			expectedBody: func() interface{} {
				return []AccountResponse{}
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "should return 200 and an slice of accounts",
			repoSetup: func() mocks.StubAccountRepository {
				return mocks.StubAccountRepository{
					Accounts: []entities.Account{acc},
				}
			},
			expectedBody: func() interface{} {
				return []AccountResponse{convertAccountToAccountResponse(acc)}
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "should return 500 and error message if something went wrong",
			repoSetup: func() mocks.StubAccountRepository {
				return mocks.StubAccountRepository{Err: errors.New("failed to list accounts")}
			},
			expectedBody: func() interface{} {
				return responses.ErrorResponse{Message: "Internal Error."}
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.repoSetup()
			accUseCase := usecase.NewAccountUseCase(&repo, nil)
			r := mux.NewRouter()
			handler := NewHandler(r, accUseCase)

			request := fakes.FakeRequest(http.MethodGet, "/accounts", nil)
			response := httptest.NewRecorder()

			http.HandlerFunc(handler.ListAccounts).ServeHTTP(response, request)

			var result []AccountResponse
			json.NewDecoder(request.Body).Decode(&result)

			assert.Equal(t, tt.expectedCode, response.Code)
			assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
			assert.Equal(t, tt.expectedBody(), result)
		})
	}

	// t.Run("should return 200 and an empty slice of accounts", func(t *testing.T) {
	// 	repo := mocks.StubAccountRepository{}
	// 	accUseCase := usecase.NewAccountUseCase(&repo, nil)
	// 	handler := NewHandler(r, accUseCase)

	// 	request, _ := http.NewRequest(http.MethodGet, "/accounts", nil)
	// 	response := httptest.NewRecorder()

	// 	http.HandlerFunc(handler.ListAccounts).ServeHTTP(response, request)

	// 	var accounts []entities.Account
	// 	json.NewDecoder(response.Body).Decode(&accounts)

	// 	assert.Equal(t, response.Code, http.StatusOK)
	// 	assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
	// 	assert.Empty(t, accounts)
	// })

	// t.Run("should return 200 and an slice of accounts", func(t *testing.T) {
	// 	acc := entities.NewAccount("Pedro", "123.456.789-00", "12345678")
	// 	repo := mocks.StubAccountRepository{Accounts: []entities.Account{acc}}
	// 	accUseCase := usecase.NewAccountUseCase(&repo, nil)
	// 	handler := NewHandler(r, accUseCase)

	// 	request := fakes.FakeRequest(http.MethodGet, "/accounts", nil)
	// 	response := httptest.NewRecorder()

	// 	http.HandlerFunc(handler.ListAccounts).ServeHTTP(response, request)

	// 	expected := []AccountResponse{convertAccountToAccountResponse(acc)}
	// 	var accounts []AccountResponse
	// 	json.NewDecoder(response.Body).Decode(&accounts)

	// 	assert.Equal(t, response.Code, http.StatusOK)
	// 	assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
	// 	assert.Equal(t, expected, accounts)
	// })

	// t.Run("should return 500 and error message if something went wrong", func(t *testing.T) {
	// 	repo := mocks.StubAccountRepository{Err: errors.New("failed to list accounts")}
	// 	accUseCase := usecase.NewAccountUseCase(&repo, nil)
	// 	handler := NewHandler(r, accUseCase)

	// 	request := fakes.FakeRequest(http.MethodGet, "/accounts", nil)
	// 	response := httptest.NewRecorder()

	// 	http.HandlerFunc(handler.ListAccounts).ServeHTTP(response, request)

	// 	expected := responses.ErrorResponse{Message: "Internal Error."}
	// 	var accounts responses.ErrorResponse
	// 	json.NewDecoder(response.Body).Decode(&accounts)

	// 	assert.Equal(t, http.StatusInternalServerError, response.Code)
	// 	assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
	// 	assert.Equal(t, expected, accounts)
	// })
}
