package account

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/thalissonfelipe/banking/pkg/domain/account/usecase"
	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/gateways/http/responses"
	"github.com/thalissonfelipe/banking/pkg/tests/mocks"
)

func TestGetAccountBalance(t *testing.T) {
	r := mux.NewRouter()

	t.Run("should return status 200 and a balance equal to 0", func(t *testing.T) {
		acc := entities.NewAccount("Pedro", "123.456.789-00", "12345678")
		repo := mocks.StubAccountRepository{Accounts: []entities.Account{acc}}
		accUseCase := usecase.NewAccountUseCase(&repo, nil)
		handler := NewHandler(r, accUseCase)

		requestURI := fmt.Sprintf("/accounts/%s/balance", acc.ID)
		request := httptest.NewRequest(http.MethodGet, requestURI, nil)
		request = mux.SetURLVars(request, map[string]string{
			"id": acc.ID,
		})
		response := httptest.NewRecorder()

		http.HandlerFunc(handler.GetAccountBalance).ServeHTTP(response, request)

		expected := 0
		var balance BalanceResponse
		json.NewDecoder(response.Body).Decode(&balance)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
		assert.Equal(t, expected, balance.Balance)
	})

	t.Run("should return status 200 and a balance equal to 100", func(t *testing.T) {
		acc := entities.NewAccount("Pedro", "123.456.789-00", "12345678")
		acc.Balance = 100
		repo := mocks.StubAccountRepository{Accounts: []entities.Account{acc}}
		accUseCase := usecase.NewAccountUseCase(&repo, nil)
		handler := NewHandler(r, accUseCase)

		requestURI := fmt.Sprintf("/accounts/%s/balance", acc.ID)
		request := httptest.NewRequest(http.MethodGet, requestURI, nil)
		request = mux.SetURLVars(request, map[string]string{
			"id": acc.ID,
		})
		response := httptest.NewRecorder()

		http.HandlerFunc(handler.GetAccountBalance).ServeHTTP(response, request)

		expected := 100
		var balance BalanceResponse
		json.NewDecoder(response.Body).Decode(&balance)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
		assert.Equal(t, expected, balance.Balance)
	})

	t.Run("should return status 404 if account does not exist", func(t *testing.T) {
		repo := mocks.StubAccountRepository{}
		accUseCase := usecase.NewAccountUseCase(&repo, nil)
		handler := NewHandler(r, accUseCase)

		requestURI := fmt.Sprintf("/accounts/%s/balance", "unknown-id")
		request := httptest.NewRequest(http.MethodGet, requestURI, nil)
		request = mux.SetURLVars(request, map[string]string{
			"id": "unknown-id",
		})
		response := httptest.NewRecorder()

		http.HandlerFunc(handler.GetAccountBalance).ServeHTTP(response, request)

		expected := responses.ErrorResponse{Message: "Account not found."}
		var result responses.ErrorResponse
		json.NewDecoder(response.Body).Decode(&result)

		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
		assert.Equal(t, expected, result)
	})

	t.Run("should return status 500 if something went wrong on usecase", func(t *testing.T) {
		acc := entities.NewAccount("Pedro", "123.456.789-00", "12345678")
		repo := mocks.StubAccountRepository{Accounts: []entities.Account{acc}, Err: errors.New("usecase fails")}
		accUseCase := usecase.NewAccountUseCase(&repo, nil)
		handler := NewHandler(r, accUseCase)

		requestURI := fmt.Sprintf("/accounts/%s/balance", acc.ID)
		request := httptest.NewRequest(http.MethodGet, requestURI, nil)
		request = mux.SetURLVars(request, map[string]string{
			"id": acc.ID,
		})
		response := httptest.NewRecorder()

		http.HandlerFunc(handler.GetAccountBalance).ServeHTTP(response, request)

		expected := responses.ErrorResponse{Message: "Internal Error."}
		var result responses.ErrorResponse
		json.NewDecoder(response.Body).Decode(&result)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
		assert.Equal(t, expected, result)
	})
}
