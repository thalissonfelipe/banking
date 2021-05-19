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
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
	"github.com/thalissonfelipe/banking/pkg/gateways/http/responses"
	"github.com/thalissonfelipe/banking/pkg/tests/mocks"
)

func TestListAccounts(t *testing.T) {
	r := mux.NewRouter()

	t.Run("should return 200 and an empty slice of accounts", func(t *testing.T) {
		repo := mocks.StubAccountRepository{}
		accUseCase := usecase.NewAccountUseCase(&repo, nil)
		handler := NewHandler(r, accUseCase)

		request, _ := http.NewRequest(http.MethodGet, "/accounts", nil)
		response := httptest.NewRecorder()

		http.HandlerFunc(handler.ListAccounts).ServeHTTP(response, request)

		var accounts []entities.Account
		json.NewDecoder(response.Body).Decode(&accounts)

		assert.Equal(t, response.Code, http.StatusOK)
		assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
		assert.Empty(t, accounts)
	})

	t.Run("should return 200 and an slice of accounts", func(t *testing.T) {
		acc := entities.NewAccount("Pedro", vos.NewCPF("123.456.789-00"), "12345678")
		repo := mocks.StubAccountRepository{Accounts: []entities.Account{acc}}
		accUseCase := usecase.NewAccountUseCase(&repo, nil)
		handler := NewHandler(r, accUseCase)

		request, _ := http.NewRequest(http.MethodGet, "/accounts", nil)
		response := httptest.NewRecorder()

		http.HandlerFunc(handler.ListAccounts).ServeHTTP(response, request)

		expected := []AccountResponse{convertAccountToAccountResponse(acc)}
		var accounts []AccountResponse
		json.NewDecoder(response.Body).Decode(&accounts)

		assert.Equal(t, response.Code, http.StatusOK)
		assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
		assert.Equal(t, expected, accounts)
	})

	t.Run("should return 500 and error message if something went wrong", func(t *testing.T) {
		repo := mocks.StubAccountRepository{Err: errors.New("failed to list accounts")}
		accUseCase := usecase.NewAccountUseCase(&repo, nil)
		handler := NewHandler(r, accUseCase)

		request, _ := http.NewRequest(http.MethodGet, "/accounts", nil)
		response := httptest.NewRecorder()

		http.HandlerFunc(handler.ListAccounts).ServeHTTP(response, request)

		expected := responses.ErrorResponse{Message: "Internal Error."}
		var accounts responses.ErrorResponse
		json.NewDecoder(response.Body).Decode(&accounts)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
		assert.Equal(t, expected, accounts)
	})
}
