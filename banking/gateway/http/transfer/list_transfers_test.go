package transfer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/domain/usecases"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
	"github.com/thalissonfelipe/banking/banking/gateway/http/rest"
	"github.com/thalissonfelipe/banking/banking/gateway/http/transfer/schema"
	"github.com/thalissonfelipe/banking/banking/gateway/jwt"
)

func TestTransferHandler_ListTransfers(t *testing.T) {
	t.Parallel()

	tr, err := entity.NewTransfer(vos.NewAccountID(), vos.NewAccountID(), 50, 100)
	require.NoError(t, err)

	transfers := []entity.Transfer{tr}

	testCases := []struct {
		name     string
		usecase  usecases.Transfer
		wantBody interface{}
		wantCode int
	}{
		{
			name: "should return a list of transfers",
			usecase: &UsecaseMock{
				ListTransfersFunc: func(context.Context, vos.AccountID) ([]entity.Transfer, error) {
					return transfers, nil
				},
			},
			wantBody: schema.MapToListTransfersResponse(transfers),
			wantCode: http.StatusOK,
		},
		{
			name: "should return an error if usecase fails",
			usecase: &UsecaseMock{
				ListTransfersFunc: func(context.Context, vos.AccountID) ([]entity.Transfer, error) {
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

			token, err := jwt.NewToken(transfers[0].AccountOriginID.String())
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, "/transfers", nil)
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
			rec := httptest.NewRecorder()

			rest.Wrap(handler.ListTransfers).ServeHTTP(rec, req)

			want, err := json.Marshal(tt.wantBody)
			require.NoError(t, err)

			assert.Equal(t, tt.wantCode, rec.Code)
			assert.JSONEq(t, string(want), rec.Body.String())
			assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
		})
	}
}
