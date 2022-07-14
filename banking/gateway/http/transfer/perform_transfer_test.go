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
		name        string
		usecase     usecases.Transfer
		accOriginID vos.AccountID
		body        interface{}
		wantBody    interface{}
		wantCode    int
	}{
		{
			name: "should perform a transfer successfully",
			usecase: &UsecaseMock{
				PerformTransferFunc: func(context.Context, usecases.PerformTransferInput) error {
					return nil
				},
			},
			accOriginID: accOrigin.ID,
			body: schema.PerformTransferInput{
				AccountDestinationID: accDest.ID.String(),
				Amount:               100,
			},
			wantBody: schema.PerformTransferResponse{
				AccountOriginID:      accOrigin.ID.String(),
				AccountDestinationID: accDest.ID.String(),
				Amount:               100,
			},
			wantCode: http.StatusCreated,
		},
		{
			name:        "should return status 400 if account dest id was not provided",
			usecase:     &UsecaseMock{},
			accOriginID: accOrigin.ID,
			body:        schema.PerformTransferInput{Amount: 100},
			wantBody: rest.Error{
				Error: "invalid request body",
				Details: []rest.ErrorDetail{
					{
						Location: "body.account_destination_id",
						Message:  "missing parameter",
					},
				},
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name:        "should return status 400 if an invalid json was provided",
			usecase:     &UsecaseMock{},
			accOriginID: accOrigin.ID,
			body: map[string]interface{}{
				"amount": "100",
			},
			wantBody: rest.Error{Error: "invalid request body"},
			wantCode: http.StatusBadRequest,
		},
		{
			name:        "should return status 400 if account destination id is invalid",
			usecase:     &UsecaseMock{},
			accOriginID: accOrigin.ID,
			body: schema.PerformTransferInput{
				AccountDestinationID: "invalid",
				Amount:               100,
			},
			wantBody: rest.Error{
				Error: "invalid request body",
				Details: []rest.ErrorDetail{
					{
						Location: "body.account_destination_id",
						Message:  "invalid uuid",
					},
				},
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "should return status 404 if acc origin does not exist",
			usecase: &UsecaseMock{
				PerformTransferFunc: func(context.Context, usecases.PerformTransferInput) error {
					return entity.ErrAccountNotFound
				},
			},
			accOriginID: accOrigin.ID,
			body: schema.PerformTransferInput{
				AccountDestinationID: accDest.ID.String(),
				Amount:               100,
			},
			wantBody: rest.Error{Error: "account origin not found"},
			wantCode: http.StatusNotFound,
		},
		{
			name: "should return status 404 if acc dest does not exist",
			usecase: &UsecaseMock{
				PerformTransferFunc: func(context.Context, usecases.PerformTransferInput) error {
					return entity.ErrAccountDestinationNotFound
				},
			},
			accOriginID: accOrigin.ID,
			body: schema.PerformTransferInput{
				AccountDestinationID: accDest.ID.String(),
				Amount:               100,
			},
			wantBody: rest.Error{Error: "account destination not found"},
			wantCode: http.StatusNotFound,
		},
		{
			name: "should return status 400 if acc origin has insufficient funds",
			usecase: &UsecaseMock{
				PerformTransferFunc: func(context.Context, usecases.PerformTransferInput) error {
					return entity.ErrInsufficientFunds
				},
			},
			accOriginID: accOrigin.ID,
			body: schema.PerformTransferInput{
				AccountDestinationID: accDest.ID.String(),
				Amount:               100,
			},
			wantBody: rest.Error{Error: "insufficient funds"},
			wantCode: http.StatusBadRequest,
		},
		{
			name:        "should return status 400 if acc dest id is the same as acc origin id",
			usecase:     &UsecaseMock{},
			accOriginID: accOrigin.ID,
			body: schema.PerformTransferInput{
				AccountDestinationID: accOrigin.ID.String(),
				Amount:               100,
			},
			wantBody: rest.Error{Error: "account origin id cannot be equal to destination id"},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "should return status 500 if usecase fails",
			usecase: &UsecaseMock{
				PerformTransferFunc: func(context.Context, usecases.PerformTransferInput) error {
					return assert.AnError
				},
			},
			accOriginID: accOrigin.ID,
			body: schema.PerformTransferInput{
				AccountDestinationID: accDest.ID.String(),
				Amount:               100,
			},
			wantBody: rest.Error{Error: "internal server error"},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandler(tt.usecase)

			request := fakes.FakeRequest(http.MethodPost, "/transfers", tt.body)
			token, _ := jwt.NewToken(tt.accOriginID.String())
			bearerToken := fmt.Sprintf("Bearer %s", token)
			request.Header.Add("Authorization", bearerToken)
			response := httptest.NewRecorder()

			rest.Wrap(handler.PerformTransfer).ServeHTTP(response, request)

			want, err := json.Marshal(tt.wantBody)
			require.NoError(t, err)

			assert.Equal(t, tt.wantCode, response.Code)
			assert.JSONEq(t, string(want), response.Body.String())
			assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
		})
	}
}
