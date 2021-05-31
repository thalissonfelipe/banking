package transfer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/transfer/usecase"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
	"github.com/thalissonfelipe/banking/pkg/gateways/http/responses"
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
		repo         *mocks.StubTransferRepository
		accUsecase   *mocks.StubAccountUsecase
		decoder      tests.Decoder
		accOriginID  vos.ID
		body         interface{}
		expectedBody interface{}
		expectedCode int
	}{
		{
			name:        "should return status 400 if account dest id was not provided",
			repo:        &mocks.StubTransferRepository{},
			accUsecase:  &mocks.StubAccountUsecase{},
			decoder:     tests.ErrorMessageDecoder{},
			accOriginID: accOrigin.ID,
			body:        transferRequest{Amount: 100},
			expectedBody: responses.ErrorResponse{
				Message: "missing account destination id parameter",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "should return status 400 if amount was not provided",
			repo:         &mocks.StubTransferRepository{},
			accUsecase:   &mocks.StubAccountUsecase{},
			decoder:      tests.ErrorMessageDecoder{},
			accOriginID:  accOrigin.ID,
			body:         transferRequest{AccountDestinationID: accDest.ID.String()},
			expectedBody: responses.ErrorResponse{Message: "missing amount parameter"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "should return status 400 an invalid json was provided",
			repo:        &mocks.StubTransferRepository{},
			accUsecase:  &mocks.StubAccountUsecase{},
			decoder:     tests.ErrorMessageDecoder{},
			accOriginID: accOrigin.ID,
			body: map[string]interface{}{
				"amount": "100",
			},
			expectedBody: responses.ErrorResponse{Message: "invalid json"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "should return status 404 if acc origin does not exist",
			repo:        &mocks.StubTransferRepository{},
			accUsecase:  &mocks.StubAccountUsecase{},
			decoder:     tests.ErrorMessageDecoder{},
			accOriginID: accOrigin.ID,
			body: transferRequest{
				AccountDestinationID: accDest.ID.String(),
				Amount:               100,
			},
			expectedBody: responses.ErrorResponse{Message: "account origin does not exist"},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "should return status 404 if acc dest does not exist",
			repo: &mocks.StubTransferRepository{},
			accUsecase: &mocks.StubAccountUsecase{
				Accounts: []entities.Account{accOrigin},
			},
			decoder:     tests.ErrorMessageDecoder{},
			accOriginID: accOrigin.ID,
			body: transferRequest{
				AccountDestinationID: accDest.ID.String(),
				Amount:               100,
			},
			expectedBody: responses.ErrorResponse{Message: "account destination does not exist"},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "should return status 400 if acc origin has insufficient funds",
			repo: &mocks.StubTransferRepository{},
			accUsecase: &mocks.StubAccountUsecase{
				Accounts: []entities.Account{accOrigin, accDest},
			},
			decoder:     tests.ErrorMessageDecoder{},
			accOriginID: accOrigin.ID,
			body: transferRequest{
				AccountDestinationID: accDest.ID.String(),
				Amount:               100,
			},
			expectedBody: responses.ErrorResponse{Message: "insufficient funds"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "should return status 500 if usecase fails",
			repo: &mocks.StubTransferRepository{},
			accUsecase: &mocks.StubAccountUsecase{
				Err: errors.New("usecase fails to fetch"),
			},
			decoder:     tests.ErrorMessageDecoder{},
			accOriginID: accOrigin.ID,
			body: transferRequest{
				AccountDestinationID: accDest.ID.String(),
				Amount:               100,
			},
			expectedBody: responses.ErrorResponse{Message: "internal server error"},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:        "should return status 400 if accDestID is the account origin id",
			repo:        &mocks.StubTransferRepository{},
			accUsecase:  &mocks.StubAccountUsecase{},
			decoder:     tests.ErrorMessageDecoder{},
			accOriginID: accOrigin.ID,
			body: transferRequest{
				AccountDestinationID: accOrigin.ID.String(),
				Amount:               100,
			},
			expectedBody: responses.ErrorResponse{
				Message: "account destination cannot be the account origin id",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "should transfer successfully",
			repo: &mocks.StubTransferRepository{},
			accUsecase: &mocks.StubAccountUsecase{
				Accounts: []entities.Account{accOriginWithBalance, accDest},
			},
			decoder:     createdTransferDecoder{},
			accOriginID: accOriginWithBalance.ID,
			body: transferRequest{
				AccountDestinationID: accDest.ID.String(),
				Amount:               100,
			},
			expectedBody: transferCreatedResponse{
				AccountOriginID:      accOriginWithBalance.ID.String(),
				AccountDestinationID: accDest.ID.String(),
				Amount:               100,
			},
			expectedCode: http.StatusOK,
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

			result := tt.decoder.Decode(response.Body)

			assert.Equal(t, tt.expectedBody, result)
			assert.Equal(t, tt.expectedCode, response.Code)
			assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
		})
	}
}

type createdTransferDecoder struct{}

func (createdTransferDecoder) Decode(body *bytes.Buffer) interface{} {
	var result transferCreatedResponse
	json.NewDecoder(body).Decode(&result)
	return result
}
