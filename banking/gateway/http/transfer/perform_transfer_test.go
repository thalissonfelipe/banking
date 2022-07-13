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

	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/domain/transfer"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
	"github.com/thalissonfelipe/banking/banking/gateway/http/rest"
	"github.com/thalissonfelipe/banking/banking/gateway/http/transfer/schema"
	"github.com/thalissonfelipe/banking/banking/services/auth"
	"github.com/thalissonfelipe/banking/banking/tests"
	"github.com/thalissonfelipe/banking/banking/tests/fakes"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
)

func TestTransferHandler_PerformTransfer(t *testing.T) {
	cpf := testdata.GetValidCPF()
	secret := testdata.GetValidSecret()

	accOrigin, err := entity.NewAccount("origin", cpf.String(), secret.String())
	require.NoError(t, err)

	accDest, err := entity.NewAccount("dest", cpf.String(), secret.String())
	require.NoError(t, err)

	testCases := []struct {
		name         string
		usecase      transfer.Usecase
		decoder      tests.Decoder
		accOriginID  vos.AccountID
		body         interface{}
		expectedBody interface{}
		expectedCode int
	}{
		{
			name: "should perform a transfer successfully",
			usecase: &UsecaseMock{
				PerformTransferFunc: func(context.Context, transfer.PerformTransferInput) error {
					return nil
				},
			},
			decoder:     createdTransferDecoder{},
			accOriginID: accOrigin.ID,
			body: schema.PerformTransferInput{
				AccountDestinationID: accDest.ID.String(),
				Amount:               100,
			},
			expectedBody: schema.PerformTransferResponse{
				AccountOriginID:      accOrigin.ID.String(),
				AccountDestinationID: accDest.ID.String(),
				Amount:               100,
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:        "should return status 400 if account dest id was not provided",
			usecase:     &UsecaseMock{},
			decoder:     tests.ErrorMessageDecoder{},
			accOriginID: accOrigin.ID,
			body:        schema.PerformTransferInput{Amount: 100},
			expectedBody: rest.ErrorResponse{
				Message: "missing account destination id parameter",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "should return status 400 if amount was not provided",
			usecase:      &UsecaseMock{},
			decoder:      tests.ErrorMessageDecoder{},
			accOriginID:  accOrigin.ID,
			body:         schema.PerformTransferInput{AccountDestinationID: accDest.ID.String()},
			expectedBody: rest.ErrorResponse{Message: "missing amount parameter"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "should return status 400 an invalid json was provided",
			usecase:     &UsecaseMock{},
			decoder:     tests.ErrorMessageDecoder{},
			accOriginID: accOrigin.ID,
			body: map[string]interface{}{
				"amount": "100",
			},
			expectedBody: rest.ErrorResponse{Message: "invalid json"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "should return status 404 if acc origin does not exist",
			usecase: &UsecaseMock{
				PerformTransferFunc: func(context.Context, transfer.PerformTransferInput) error {
					return entity.ErrAccountNotFound
				},
			},
			decoder:     tests.ErrorMessageDecoder{},
			accOriginID: accOrigin.ID,
			body: schema.PerformTransferInput{
				AccountDestinationID: accDest.ID.String(),
				Amount:               100,
			},
			expectedBody: rest.ErrorResponse{Message: "account origin does not exist"},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "should return status 404 if acc dest does not exist",
			usecase: &UsecaseMock{
				PerformTransferFunc: func(context.Context, transfer.PerformTransferInput) error {
					return entity.ErrAccountDestinationNotFound
				},
			},
			decoder:     tests.ErrorMessageDecoder{},
			accOriginID: accOrigin.ID,
			body: schema.PerformTransferInput{
				AccountDestinationID: accDest.ID.String(),
				Amount:               100,
			},
			expectedBody: rest.ErrorResponse{Message: "account destination does not exist"},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "should return status 400 if acc origin has insufficient funds",
			usecase: &UsecaseMock{
				PerformTransferFunc: func(context.Context, transfer.PerformTransferInput) error {
					return entity.ErrInsufficientFunds
				},
			},
			decoder:     tests.ErrorMessageDecoder{},
			accOriginID: accOrigin.ID,
			body: schema.PerformTransferInput{
				AccountDestinationID: accDest.ID.String(),
				Amount:               100,
			},
			expectedBody: rest.ErrorResponse{Message: "insufficient funds"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "should return status 400 if accDestID is the same as account origin id",
			usecase:     &UsecaseMock{},
			decoder:     tests.ErrorMessageDecoder{},
			accOriginID: accOrigin.ID,
			body: schema.PerformTransferInput{
				AccountDestinationID: accOrigin.ID.String(),
				Amount:               100,
			},
			expectedBody: rest.ErrorResponse{
				Message: "account destination cannot be the account origin id",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "should return status 500 if usecase fails",
			usecase: &UsecaseMock{
				PerformTransferFunc: func(context.Context, transfer.PerformTransferInput) error {
					return assert.AnError
				},
			},
			decoder:     tests.ErrorMessageDecoder{},
			accOriginID: accOrigin.ID,
			body: schema.PerformTransferInput{
				AccountDestinationID: accDest.ID.String(),
				Amount:               100,
			},
			expectedBody: rest.ErrorResponse{Message: "internal server error"},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			handler := NewHandler(r, tt.usecase)

			request := fakes.FakeRequest(http.MethodPost, "/transfers", tt.body)
			token, _ := auth.NewToken(tt.accOriginID.String())
			bearerToken := fmt.Sprintf("Bearer %s", token)
			request.Header.Add("Authorization", bearerToken)
			response := httptest.NewRecorder()

			http.HandlerFunc(handler.PerformTransfer).ServeHTTP(response, request)

			result := tt.decoder.Decode(t, response.Body)

			assert.Equal(t, tt.expectedBody, result)
			assert.Equal(t, tt.expectedCode, response.Code)
			assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
		})
	}
}

type createdTransferDecoder struct{}

func (createdTransferDecoder) Decode(t *testing.T, body *bytes.Buffer) interface{} {
	var result schema.PerformTransferResponse

	err := json.NewDecoder(body).Decode(&result)
	require.NoError(t, err)

	return result
}
