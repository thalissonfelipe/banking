package account

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func TestHandler_GetAccountBalance(t *testing.T) {
	acc := entities.NewAccount("Pedro", testdata.GetValidCPF(), testdata.GetValidSecret())

	testCases := []struct {
		name         string
		repo         func() *mocks.StubAccountRepository
		requestURI   string
		decoder      tests.Decoder
		expectedBody interface{}
		expectedCode int
	}{
		{
			name: "should return status 200 and a balance equal to 0",
			repo: func() *mocks.StubAccountRepository {
				return &mocks.StubAccountRepository{
					Accounts: []entities.Account{acc},
				}
			},
			requestURI:   fmt.Sprintf("/accounts/%s/balance", acc.ID),
			decoder:      balanceResponseDecoder{},
			expectedBody: balanceResponse{Balance: 0},
			expectedCode: http.StatusOK,
		},
		{
			name: "should return status 200 and a balance equal to 100",
			repo: func() *mocks.StubAccountRepository {
				accWithBalance := acc
				accWithBalance.Balance = 100
				return &mocks.StubAccountRepository{
					Accounts: []entities.Account{accWithBalance},
				}
			},
			requestURI:   fmt.Sprintf("/accounts/%s/balance", acc.ID),
			decoder:      balanceResponseDecoder{},
			expectedBody: balanceResponse{Balance: 100},
			expectedCode: http.StatusOK,
		},
		{
			name: "should return status 404 if account does not exist",
			repo: func() *mocks.StubAccountRepository {
				return &mocks.StubAccountRepository{}
			},
			requestURI:   fmt.Sprintf("/accounts/%s/balance", acc.ID),
			decoder:      tests.ErrorMessageDecoder{},
			expectedBody: responses.ErrorResponse{Message: "account does not exist"},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "should return status 500 if something went wrong on usecase",
			repo: func() *mocks.StubAccountRepository {
				return &mocks.StubAccountRepository{
					Err: entities.ErrInternalError,
				}
			},
			requestURI:   fmt.Sprintf("/accounts/%s/balance", acc.ID),
			decoder:      tests.ErrorMessageDecoder{},
			expectedBody: responses.ErrorResponse{Message: "internal server error"},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			r := mux.NewRouter()
			usecase := usecase.NewAccountUsecase(tt.repo(), nil)
			handler := NewHandler(r, usecase)

			request := fakes.FakeRequest(http.MethodGet, tt.requestURI, nil)
			request = mux.SetURLVars(request, map[string]string{
				"id": acc.ID.String(),
			})
			response := httptest.NewRecorder()

			http.HandlerFunc(handler.GetAccountBalance).ServeHTTP(response, request)

			result := tt.decoder.Decode(response.Body)

			assert.Equal(t, tt.expectedBody, result)
			assert.Equal(t, tt.expectedCode, response.Code)
			assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
		})
	}
}

type balanceResponseDecoder struct{}

func (balanceResponseDecoder) Decode(body *bytes.Buffer) interface{} {
	var result balanceResponse
	json.NewDecoder(body).Decode(&result)
	return result
}
