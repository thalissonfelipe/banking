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
	"github.com/thalissonfelipe/banking/pkg/tests/testdata"
)

func TestHandler_ListTransfers(t *testing.T) {
	accOrigin := entities.NewAccount("Pedro", testdata.GetValidCPF(), vos.NewSecret("12345678"))
	accDest := entities.NewAccount("Maria", testdata.GetValidCPF(), vos.NewSecret("12345678"))
	transfer := entities.NewTransfer(accOrigin.ID, accDest.ID, 100)

	testCases := []struct {
		name         string
		repo         *mocks.StubTransferRepository
		decoder      tests.Decoder
		expectedBody interface{}
		expectedCode int
	}{
		{
			name:         "should return a empty list of transfers",
			repo:         &mocks.StubTransferRepository{},
			decoder:      listTransfersDecoder{},
			expectedBody: []transferResponse{},
			expectedCode: http.StatusOK,
		},
		{
			name: "should return a list of transfers",
			repo: &mocks.StubTransferRepository{
				Transfers: []entities.Transfer{transfer},
			},
			decoder:      listTransfersDecoder{},
			expectedBody: []transferResponse{convertTransferToTransferResponse(transfer)},
			expectedCode: http.StatusOK,
		},
		{
			name: "should return an error if usecase fails",
			repo: &mocks.StubTransferRepository{
				Err: errors.New("usecase fails to fetch transfers"),
			},
			decoder:      tests.ErrorMessageDecoder{},
			expectedBody: responses.ErrorResponse{Message: "internal server error"},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			r := mux.NewRouter()
			accUsecase := &mocks.StubAccountUsecase{}
			trUsecase := usecase.NewTransferUsecase(tt.repo, accUsecase)
			handler := NewHandler(r, trUsecase)

			request := fakes.FakeRequest(http.MethodGet, "/transfers", nil)
			token, _ := auth.NewToken(accOrigin.ID.String())
			bearerToken := fmt.Sprintf("Bearer %s", token)
			request.Header.Add("Authorization", bearerToken)
			response := httptest.NewRecorder()

			http.HandlerFunc(handler.ListTransfers).ServeHTTP(response, request)

			result := tt.decoder.Decode(response.Body)

			assert.Equal(t, tt.expectedBody, result)
			assert.Equal(t, tt.expectedCode, response.Code)
			assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
		})
	}
}

type listTransfersDecoder struct{}

func (listTransfersDecoder) Decode(body *bytes.Buffer) interface{} {
	var result []transferResponse
	json.NewDecoder(body).Decode(&result)
	return result
}
