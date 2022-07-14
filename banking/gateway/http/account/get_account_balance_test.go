package account

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/domain/usecases"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
	"github.com/thalissonfelipe/banking/banking/gateway/http/account/schema"
	"github.com/thalissonfelipe/banking/banking/gateway/http/rest"
	"github.com/thalissonfelipe/banking/banking/tests/fakes"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
)

func TestAccountHandler_GetAccountBalance(t *testing.T) {
	acc, err := entity.NewAccount("name", testdata.GetValidCPF().String(), testdata.GetValidSecret().String())
	require.NoError(t, err)

	testCases := []struct {
		name      string
		usecase   usecases.Account
		accountID string
		wantBody  interface{}
		wantCode  int
	}{
		{
			name: "should return account balance successfully",
			usecase: &UsecaseMock{
				GetAccountBalanceByIDFunc: func(context.Context, vos.AccountID) (int, error) {
					return 100, nil
				},
			},
			accountID: acc.ID.String(),
			wantBody:  schema.BalanceResponse{Balance: 100},
			wantCode:  http.StatusOK,
		},
		{
			name:      "should return status 400 if account id is invalid",
			usecase:   &UsecaseMock{},
			accountID: "invalid",
			wantBody: rest.Error{
				Error: "invalid path parameters",
				Details: []rest.ErrorDetail{
					{
						Location: "path.accountID",
						Message:  "invalid uuid",
					},
				},
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "should return status 404 if account does not exist",
			usecase: &UsecaseMock{
				GetAccountBalanceByIDFunc: func(context.Context, vos.AccountID) (int, error) {
					return 0, entity.ErrAccountNotFound
				},
			},
			accountID: acc.ID.String(),
			wantBody:  rest.Error{Error: "account not found"},
			wantCode:  http.StatusNotFound,
		},
		{
			name: "should return status 500 if usecase fails",
			usecase: &UsecaseMock{
				GetAccountBalanceByIDFunc: func(context.Context, vos.AccountID) (int, error) {
					return 0, assert.AnError
				},
			},
			accountID: acc.ID.String(),
			wantBody:  rest.Error{Error: "internal server error"},
			wantCode:  http.StatusInternalServerError,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandler(tt.usecase)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("accountID", tt.accountID)

			requestURI := fmt.Sprintf("/accounts/%s/balance", tt.accountID)

			request := fakes.FakeRequest(http.MethodGet, requestURI, nil)
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

			response := httptest.NewRecorder()

			rest.Wrap(handler.GetAccountBalance).ServeHTTP(response, request)

			want, err := json.Marshal(tt.wantBody)
			require.NoError(t, err)

			assert.Equal(t, tt.wantCode, response.Code)
			assert.JSONEq(t, string(want), response.Body.String())
			assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
		})
	}
}
