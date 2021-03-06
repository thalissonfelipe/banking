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
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
)

func TestAccountHandler_ListAccounts(t *testing.T) {
	t.Parallel()

	acc, err := entity.NewAccount("name", testdata.CPF().String(), testdata.Secret().String())
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
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := NewHandler(tt.usecase)

			req := httptest.NewRequest(http.MethodGet, "/accounts", nil)
			rec := httptest.NewRecorder()

			rest.Wrap(handler.ListAccounts).ServeHTTP(rec, req)

			want, err := json.Marshal(tt.wantBody)
			require.NoError(t, err)

			assert.Equal(t, tt.wantCode, rec.Code)
			assert.JSONEq(t, string(want), rec.Body.String())
			assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
		})
	}
}
