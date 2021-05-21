package transfer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/transfer/usecase"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
	"github.com/thalissonfelipe/banking/pkg/gateways/http/responses"
	"github.com/thalissonfelipe/banking/pkg/services/auth"
	"github.com/thalissonfelipe/banking/pkg/tests"
	"github.com/thalissonfelipe/banking/pkg/tests/fakes"
	"github.com/thalissonfelipe/banking/pkg/tests/mocks"
)

func TestCreateTransfer(t *testing.T) {
	accOrigin := entities.NewAccount("Pedro", vos.NewCPF("123.456.789-00"), vos.NewSecret("12345678"))
	accDest := entities.NewAccount("Maria", vos.NewCPF("123.456.789-11"), vos.NewSecret("12345678"))
	accOriginWithBalance := entities.NewAccount("João", vos.NewCPF("123.456.789-22"), vos.NewSecret("12345678"))
	accOriginWithBalance.Balance = 200

	testCases := []struct {
		name         string
		repo         *mocks.StubTransferRepository
		accUsecase   *mocks.StubAccountUseCase
		decoder      tests.Decoder
		accOriginID  string
		body         interface{}
		expectedBody interface{}
		expectedCode int
	}{
		{
			name:        "should return status 400 if account dest id was not provided",
			repo:        &mocks.StubTransferRepository{},
			accUsecase:  &mocks.StubAccountUseCase{},
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
			accUsecase:   &mocks.StubAccountUseCase{},
			decoder:      tests.ErrorMessageDecoder{},
			accOriginID:  accOrigin.ID,
			body:         transferRequest{AccountDestinationID: accDest.ID},
			expectedBody: responses.ErrorResponse{Message: "missing amount parameter"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "should return status 400 an invalid json was provided",
			repo:        &mocks.StubTransferRepository{},
			accUsecase:  &mocks.StubAccountUseCase{},
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
			accUsecase:  &mocks.StubAccountUseCase{},
			decoder:     tests.ErrorMessageDecoder{},
			accOriginID: accOrigin.ID,
			body: transferRequest{
				AccountDestinationID: accDest.ID,
				Amount:               100,
			},
			expectedBody: responses.ErrorResponse{Message: "account origin does not exist"},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "should return status 404 if acc dest does not exist",
			repo: &mocks.StubTransferRepository{},
			accUsecase: &mocks.StubAccountUseCase{
				Accounts: []entities.Account{accOrigin},
			},
			decoder:     tests.ErrorMessageDecoder{},
			accOriginID: accOrigin.ID,
			body: transferRequest{
				AccountDestinationID: accDest.ID,
				Amount:               100,
			},
			expectedBody: responses.ErrorResponse{Message: "account destination does not exist"},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "should return status 400 if acc origin has insufficient funds",
			repo: &mocks.StubTransferRepository{},
			accUsecase: &mocks.StubAccountUseCase{
				Accounts: []entities.Account{accOrigin, accDest},
			},
			decoder:     tests.ErrorMessageDecoder{},
			accOriginID: accOrigin.ID,
			body: transferRequest{
				AccountDestinationID: accDest.ID,
				Amount:               100,
			},
			expectedBody: responses.ErrorResponse{Message: "insufficient funds"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "should return status 500 if usecase fails",
			repo: &mocks.StubTransferRepository{},
			accUsecase: &mocks.StubAccountUseCase{
				Err: errors.New("usecase fails to fetch"),
			},
			decoder:     tests.ErrorMessageDecoder{},
			accOriginID: accOrigin.ID,
			body: transferRequest{
				AccountDestinationID: accDest.ID,
				Amount:               100,
			},
			expectedBody: responses.ErrorResponse{Message: "internal server error"},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:        "should return status 400 if accDestID is the account origin id",
			repo:        &mocks.StubTransferRepository{},
			accUsecase:  &mocks.StubAccountUseCase{},
			decoder:     tests.ErrorMessageDecoder{},
			accOriginID: accOrigin.ID,
			body: transferRequest{
				AccountDestinationID: accOrigin.ID,
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
			accUsecase: &mocks.StubAccountUseCase{
				Accounts: []entities.Account{accOriginWithBalance, accDest},
			},
			decoder:     createdTransferDecoder{},
			accOriginID: accOriginWithBalance.ID,
			body: transferRequest{
				AccountDestinationID: accDest.ID,
				Amount:               100,
			},
			expectedBody: transferCreatedResponse{
				AccountOriginID:      accOriginWithBalance.ID,
				AccountDestinationID: accDest.ID,
				Amount:               100,
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			r := mux.NewRouter()
			usecase := usecase.NewTransfer(tt.repo, tt.accUsecase)
			handler := NewHandler(r, usecase)

			request := fakes.FakeRequest(http.MethodPost, "/transfers", tt.body)
			token, _ := auth.NewToken(tt.accOriginID)
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