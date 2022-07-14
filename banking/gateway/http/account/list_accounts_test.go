package account

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/domain/usecases"
	"github.com/thalissonfelipe/banking/banking/gateway/http/account/schema"
	"github.com/thalissonfelipe/banking/banking/gateway/http/rest"
	"github.com/thalissonfelipe/banking/banking/tests/fakes"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
)

func TestAccountHandler_ListAccounts(t *testing.T) {
	acc, err := entity.NewAccount("name", testdata.GetValidCPF().String(), testdata.GetValidSecret().String())
	require.NoError(t, err)

	accounts := []entity.Account{acc}

	testCases := []struct {
		name     string
		usecase  usecases.Account
		wantBody interface{}
		wantCode int
	}{
		{
			name: "should return a list of accounts successfully",
			usecase: &UsecaseMock{
				ListAccountsFunc: func(context.Context) ([]entity.Account, error) {
					return accounts, nil
				},
			},
			wantBody: schema.MapToListAccountsResponse(accounts),
			wantCode: http.StatusOK,
		},
		{
			name: "should return 500 if usecase fails",
			usecase: &UsecaseMock{
				ListAccountsFunc: func(context.Context) ([]entity.Account, error) {
					return nil, assert.AnError
				},
			},
			wantBody: rest.Error{Error: "internal server error"},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandler(tt.usecase)

			request := fakes.FakeRequest(http.MethodGet, "/accounts", nil)
			response := httptest.NewRecorder()

			rest.Wrap(handler.ListAccounts).ServeHTTP(response, request)

			want, err := json.Marshal(tt.wantBody)
			require.NoError(t, err)

			assert.Equal(t, tt.wantCode, response.Code)
			assert.JSONEq(t, string(want), response.Body.String())
			assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
		})
	}
}
