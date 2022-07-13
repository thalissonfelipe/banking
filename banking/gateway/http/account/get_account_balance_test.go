package account

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/banking/domain/account"
	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
	"github.com/thalissonfelipe/banking/banking/gateway/http/account/schema"
	"github.com/thalissonfelipe/banking/banking/gateway/http/rest"
	"github.com/thalissonfelipe/banking/banking/tests"
	"github.com/thalissonfelipe/banking/banking/tests/fakes"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
)

func TestAccountHandler_GetAccountBalance(t *testing.T) {
	acc, err := entity.NewAccount("name", testdata.GetValidCPF().String(), testdata.GetValidSecret().String())
	require.NoError(t, err)

	testCases := []struct {
		name         string
		usecase      account.Usecase
		requestURI   string
		decoder      tests.Decoder
		expectedBody interface{}
		expectedCode int
	}{
		{
			name: "should return account balance successfully",
			usecase: &UsecaseMock{
				GetAccountBalanceByIDFunc: func(context.Context, vos.AccountID) (int, error) {
					return 100, nil
				},
			},
			requestURI:   fmt.Sprintf("/accounts/%s/balance", acc.ID),
			decoder:      balanceResponseDecoder{},
			expectedBody: schema.BalanceResponse{Balance: 100},
			expectedCode: http.StatusOK,
		},
		{
			name: "should return status 404 if account does not exist",
			usecase: &UsecaseMock{
				GetAccountBalanceByIDFunc: func(context.Context, vos.AccountID) (int, error) {
					return 0, entity.ErrAccountNotFound
				},
			},
			requestURI:   fmt.Sprintf("/accounts/%s/balance", acc.ID),
			decoder:      tests.ErrorMessageDecoder{},
			expectedBody: rest.ErrorResponse{Message: "account does not exist"},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "should return status 500 if usecase fails",
			usecase: &UsecaseMock{
				GetAccountBalanceByIDFunc: func(context.Context, vos.AccountID) (int, error) {
					return 0, assert.AnError
				},
			},
			requestURI:   fmt.Sprintf("/accounts/%s/balance", acc.ID),
			decoder:      tests.ErrorMessageDecoder{},
			expectedBody: rest.ErrorResponse{Message: "internal server error"},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			handler := NewHandler(r, tt.usecase)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("accountID", acc.ID.String())

			request := fakes.FakeRequest(http.MethodGet, tt.requestURI, nil)
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

			response := httptest.NewRecorder()

			http.HandlerFunc(handler.GetAccountBalance).ServeHTTP(response, request)

			result := tt.decoder.Decode(t, response.Body)

			assert.Equal(t, tt.expectedBody, result)
			assert.Equal(t, tt.expectedCode, response.Code)
			assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
		})
	}
}

type balanceResponseDecoder struct{}

func (balanceResponseDecoder) Decode(t *testing.T, body *bytes.Buffer) interface{} {
	var result schema.BalanceResponse

	err := json.NewDecoder(body).Decode(&result)
	require.NoError(t, err)

	return result
}
