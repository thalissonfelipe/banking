package transfer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/transfer/usecase"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
	"github.com/thalissonfelipe/banking/pkg/gateways/http/rest"
	"github.com/thalissonfelipe/banking/pkg/gateways/http/transfer/schemes"
	"github.com/thalissonfelipe/banking/pkg/services/auth"
	"github.com/thalissonfelipe/banking/pkg/tests"
	"github.com/thalissonfelipe/banking/pkg/tests/fakes"
	"github.com/thalissonfelipe/banking/pkg/tests/mocks"
	"github.com/thalissonfelipe/banking/pkg/tests/testdata"
)

func TestHandler_CreateTransfer(t *testing.T) {
	accOrigin := entities.NewAccount("Pedro", testdata.GetValidCPF(), testdata.GetValidSecret())
	accDest := entities.NewAccount("Maria", testdata.GetValidCPF(), testdata.GetValidSecret())
	accOriginWithBalance := entities.NewAccount("João", testdata.GetValidCPF(), testdata.GetValidSecret())
	accOriginWithBalance.Balance = 200

	testCases := []struct {
		name         string
		repo         *mocks.TransferRepositoryMock
		accUsecase   *mocks.AccountUsecaseMock
		decoder      tests.Decoder
		accOriginID  vos.AccountID
		body         interface{}
		expectedBody interface{}
		expectedCode int
	}{
		{
			name:        "should return status 400 if account dest id was not provided",
			repo:        &mocks.TransferRepositoryMock{},
			accUsecase:  &mocks.AccountUsecaseMock{},
			decoder:     tests.ErrorMessageDecoder{},
			accOriginID: accOrigin.ID,
			body:        schemes.CreateTransferInput{Amount: 100},
			expectedBody: rest.ErrorResponse{
				Message: "missing account destination id parameter",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "should return status 400 if amount was not provided",
			repo:         &mocks.TransferRepositoryMock{},
			accUsecase:   &mocks.AccountUsecaseMock{},
			decoder:      tests.ErrorMessageDecoder{},
			accOriginID:  accOrigin.ID,
			body:         schemes.CreateTransferInput{AccountDestinationID: accDest.ID.String()},
			expectedBody: rest.ErrorResponse{Message: "missing amount parameter"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "should return status 400 an invalid json was provided",
			repo:        &mocks.TransferRepositoryMock{},
			accUsecase:  &mocks.AccountUsecaseMock{},
			decoder:     tests.ErrorMessageDecoder{},
			accOriginID: accOrigin.ID,
			body: map[string]interface{}{
				"amount": "100",
			},
			expectedBody: rest.ErrorResponse{Message: "invalid json"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "should return status 404 if acc origin does not exist",
			repo:        &mocks.TransferRepositoryMock{},
			accUsecase:  &mocks.AccountUsecaseMock{},
			decoder:     tests.ErrorMessageDecoder{},
			accOriginID: accOrigin.ID,
			body: schemes.CreateTransferInput{
				AccountDestinationID: accDest.ID.String(),
				Amount:               100,
			},
			expectedBody: rest.ErrorResponse{Message: "account origin does not exist"},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "should return status 404 if acc dest does not exist",
			repo: &mocks.TransferRepositoryMock{},
			accUsecase: &mocks.AccountUsecaseMock{
				Accounts: []entities.Account{accOrigin},
			},
			decoder:     tests.ErrorMessageDecoder{},
			accOriginID: accOrigin.ID,
			body: schemes.CreateTransferInput{
				AccountDestinationID: accDest.ID.String(),
				Amount:               100,
			},
			expectedBody: rest.ErrorResponse{Message: "account destination does not exist"},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "should return status 400 if acc origin has insufficient funds",
			repo: &mocks.TransferRepositoryMock{},
			accUsecase: &mocks.AccountUsecaseMock{
				Accounts: []entities.Account{accOrigin, accDest},
			},
			decoder:     tests.ErrorMessageDecoder{},
			accOriginID: accOrigin.ID,
			body: schemes.CreateTransferInput{
				AccountDestinationID: accDest.ID.String(),
				Amount:               100,
			},
			expectedBody: rest.ErrorResponse{Message: "insufficient funds"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "should return status 500 if usecase fails",
			repo:        &mocks.TransferRepositoryMock{},
			accUsecase:  &mocks.AccountUsecaseMock{Err: testdata.ErrUsecaseFails},
			decoder:     tests.ErrorMessageDecoder{},
			accOriginID: accOrigin.ID,
			body: schemes.CreateTransferInput{
				AccountDestinationID: accDest.ID.String(),
				Amount:               100,
			},
			expectedBody: rest.ErrorResponse{Message: "internal server error"},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:        "should return status 400 if accDestID is the account origin id",
			repo:        &mocks.TransferRepositoryMock{},
			accUsecase:  &mocks.AccountUsecaseMock{},
			decoder:     tests.ErrorMessageDecoder{},
			accOriginID: accOrigin.ID,
			body: schemes.CreateTransferInput{
				AccountDestinationID: accOrigin.ID.String(),
				Amount:               100,
			},
			expectedBody: rest.ErrorResponse{
				Message: "account destination cannot be the account origin id",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "should transfer successfully",
			repo: &mocks.TransferRepositoryMock{},
			accUsecase: &mocks.AccountUsecaseMock{
				Accounts: []entities.Account{accOriginWithBalance, accDest},
			},
			decoder:     createdTransferDecoder{},
			accOriginID: accOriginWithBalance.ID,
			body: schemes.CreateTransferInput{
				AccountDestinationID: accDest.ID.String(),
				Amount:               100,
			},
			expectedBody: schemes.CreateTransferResponse{
				AccountOriginID:      accOriginWithBalance.ID.String(),
				AccountDestinationID: accDest.ID.String(),
				Amount:               100,
			},
			expectedCode: http.StatusCreated,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			usecase := usecase.NewTransferUsecase(tt.repo, tt.accUsecase)
			handler := NewHandler(r, usecase)

			request := fakes.FakeRequest(http.MethodPost, "/transfers", tt.body)
			token, _ := auth.NewToken(tt.accOriginID.String())
			bearerToken := fmt.Sprintf("Bearer %s", token)
			request.Header.Add("Authorization", bearerToken)
			response := httptest.NewRecorder()

			http.HandlerFunc(handler.CreateTransfer).ServeHTTP(response, request)

			result := tt.decoder.Decode(t, response.Body)

			assert.Equal(t, tt.expectedBody, result)
			assert.Equal(t, tt.expectedCode, response.Code)
			assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
		})
	}
}

type createdTransferDecoder struct{}

func (createdTransferDecoder) Decode(t *testing.T, body *bytes.Buffer) interface{} {
	var result schemes.CreateTransferResponse

	err := json.NewDecoder(body).Decode(&result)
	require.NoError(t, err)

	return result
}
