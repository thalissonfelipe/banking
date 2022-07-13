package transfer

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

	"github.com/thalissonfelipe/banking/banking/domain/entities"
	"github.com/thalissonfelipe/banking/banking/domain/transfer"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
	"github.com/thalissonfelipe/banking/banking/gateway/http/rest"
	"github.com/thalissonfelipe/banking/banking/gateway/http/transfer/schema"
	"github.com/thalissonfelipe/banking/banking/services/auth"
	"github.com/thalissonfelipe/banking/banking/tests"
	"github.com/thalissonfelipe/banking/banking/tests/fakes"
)

func TestTransferHandler_ListTransfers(t *testing.T) {
	tr, err := entities.NewTransfer(vos.NewAccountID(), vos.NewAccountID(), 50, 100)
	require.NoError(t, err)

	transfers := []entities.Transfer{tr}

	testCases := []struct {
		name         string
		usecase      transfer.Usecase
		decoder      tests.Decoder
		expectedBody interface{}
		expectedCode int
	}{
		{
			name: "should return a list of transfers",
			usecase: &UsecaseMock{
				ListTransfersFunc: func(context.Context, vos.AccountID) ([]entities.Transfer, error) {
					return transfers, nil
				},
			},
			decoder:      listTransfersDecoder{},
			expectedBody: schema.MapToListTransfersResponse(transfers),
			expectedCode: http.StatusOK,
		},
		{
			name: "should return an error if usecase fails",
			usecase: &UsecaseMock{
				ListTransfersFunc: func(context.Context, vos.AccountID) ([]entities.Transfer, error) {
					return nil, assert.AnError
				},
			},
			decoder:      tests.ErrorMessageDecoder{},
			expectedBody: rest.ErrorResponse{Message: "internal server error"},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			handler := NewHandler(r, tt.usecase)

			token, err := auth.NewToken(transfers[0].AccountOriginID.String())
			require.NoError(t, err)

			request := fakes.FakeRequest(http.MethodGet, "/transfers", nil)
			bearerToken := fmt.Sprintf("Bearer %s", token)
			request.Header.Add("Authorization", bearerToken)

			response := httptest.NewRecorder()

			http.HandlerFunc(handler.ListTransfers).ServeHTTP(response, request)

			result := tt.decoder.Decode(t, response.Body)

			assert.Equal(t, tt.expectedBody, result)
			assert.Equal(t, tt.expectedCode, response.Code)
			assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
		})
	}
}

type listTransfersDecoder struct{}

func (listTransfersDecoder) Decode(t *testing.T, body *bytes.Buffer) interface{} {
	var result schema.ListTransfersResponse

	err := json.NewDecoder(body).Decode(&result)
	require.NoError(t, err)

	return result
}
