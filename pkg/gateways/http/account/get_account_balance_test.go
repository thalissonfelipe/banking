package account

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/thalissonfelipe/banking/pkg/domain/account/usecase"
	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

func TestGetAccountBalance(t *testing.T) {
	r := mux.NewRouter()

	t.Run("should return status 200 and a balance equal to 0", func(t *testing.T) {
		acc := entities.NewAccount("Pedro", "123.456.789-00", "12345678")
		repo := StubAccountRepository{accounts: []entities.Account{acc}}
		accUseCase := usecase.NewAccountUseCase(&repo, nil)
		handler := NewHandler(r, accUseCase)

		requestURI := fmt.Sprintf("/accounts/%s/balance", acc.ID)
		log.Println(requestURI)
		request, _ := http.NewRequest(http.MethodGet, requestURI, nil)
		response := httptest.NewRecorder()

		http.HandlerFunc(handler.GetAccountBalance).ServeHTTP(response, request)

		expected := 0
		var balance BalanceResponse
		json.NewDecoder(response.Body).Decode(&balance)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
		assert.Equal(t, expected, balance.Balance)
	})

	// t.Run("should return status 200 and a balance equal to 100", func(t *testing.T) {
	// 	acc := entities.NewAccount("Pedro", "123.456.789-00", "12345678")
	// 	acc.Balance = 100
	// 	repo := StubAccountRepository{accounts: []entities.Account{acc}}
	// 	accUseCase := usecase.NewAccountUseCase(&repo, nil)
	// 	handler := NewHandler(r, accUseCase)

	// 	requestURI := fmt.Sprintf("/accounts/%s/balance", acc.ID)
	// 	request, _ := http.NewRequest(http.MethodGet, requestURI, nil)
	// 	response := httptest.NewRecorder()

	// 	http.HandlerFunc(handler.GetAccountBalance).ServeHTTP(response, request)

	// 	expected := 100
	// 	var balance BalanceResponse
	// 	json.NewDecoder(response.Body).Decode(&balance)

	// 	assert.Equal(t, http.StatusOK, response.Code)
	// 	assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
	// 	assert.Equal(t, expected, balance.Balance)
	// })

	// t.Run("should return status 404", func(t *testing.T) {
	// 	repo := StubAccountRepository{}
	// 	accUseCase := usecase.NewAccountUseCase(&repo, nil)
	// 	handler := NewHandler(r, accUseCase)

	// 	requestURI := fmt.Sprintf("/accounts/%s/balance", "123456")
	// 	request, _ := http.NewRequest(http.MethodGet, requestURI, nil)
	// 	response := httptest.NewRecorder()

	// 	http.HandlerFunc(handler.GetAccountBalance).ServeHTTP(response, request)

	// 	assert.Equal(t, http.StatusNotFound, response.Code)
	// 	assert.Equal(t, "Account not found.", response.Body.String())
	// })
}
