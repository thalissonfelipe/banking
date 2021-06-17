package integration

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/pkg/domain/vos"
	"github.com/thalissonfelipe/banking/pkg/gateways/http/transfer/schemes"
	"github.com/thalissonfelipe/banking/pkg/tests/dockertest"
	"github.com/thalissonfelipe/banking/pkg/tests/fakes"
	"github.com/thalissonfelipe/banking/pkg/tests/testdata"
	"github.com/thalissonfelipe/banking/pkg/tests/testenv"
)

func TestIntegration_CreateTransfer(t *testing.T) {
	uri := testenv.ServerURL + "/api/v1/transfers"

	testCases := []struct {
		name           string
		requestSetup   func(t *testing.T) *http.Request
		expectedStatus int
	}{
		{
			name: "should transfer successfully",
			requestSetup: func(t *testing.T) *http.Request {
				secret := testdata.GetValidSecret()
				acc1 := createAccount(t, testdata.GetValidCPF(), secret, 100)
				acc2 := createAccount(t, testdata.GetValidCPF(), testdata.GetValidSecret(), 0)

				reqBody := schemes.CreateTransferInput{AccountDestinationID: acc2.ID.String(), Amount: 100}
				request := fakes.FakeAuthorizedRequest(t, http.MethodPost, uri, acc1.CPF.String(), secret.String(), reqBody)

				return request
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "should return 400 if account does not has sufficient funds",
			requestSetup: func(t *testing.T) *http.Request {
				secret := testdata.GetValidSecret()
				acc1 := createAccount(t, testdata.GetValidCPF(), secret, 100)
				acc2 := createAccount(t, testdata.GetValidCPF(), testdata.GetValidSecret(), 0)

				reqBody := schemes.CreateTransferInput{AccountDestinationID: acc2.ID.String(), Amount: 200}
				request := fakes.FakeAuthorizedRequest(t, http.MethodPost, uri, acc1.CPF.String(), secret.String(), reqBody)

				return request
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "should return 400 if the account_destination_id is the same as the account origin",
			requestSetup: func(t *testing.T) *http.Request {
				secret := testdata.GetValidSecret()
				accOrigin := createAccount(t, testdata.GetValidCPF(), secret, 100)

				reqBody := schemes.CreateTransferInput{AccountDestinationID: accOrigin.ID.String(), Amount: 200}
				request := fakes.FakeAuthorizedRequest(t, http.MethodPost, uri, accOrigin.CPF.String(), secret.String(), reqBody)

				return request
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "should return 404 if account origin does not exist anymore",
			requestSetup: func(t *testing.T) *http.Request {
				secret := testdata.GetValidSecret()
				accOrigin := createAccount(t, testdata.GetValidCPF(), secret, 100)

				reqBody := schemes.CreateTransferInput{AccountDestinationID: vos.NewID().String(), Amount: 200}
				request := fakes.FakeAuthorizedRequest(t, http.MethodPost, uri, accOrigin.CPF.String(), secret.String(), reqBody)

				dockertest.TruncateTables(context.Background(), testenv.DB)

				return request
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "should return 404 if account destination does not exist",
			requestSetup: func(t *testing.T) *http.Request {
				secret := testdata.GetValidSecret()
				accOrigin := createAccount(t, testdata.GetValidCPF(), secret, 100)

				reqBody := schemes.CreateTransferInput{AccountDestinationID: vos.NewID().String(), Amount: 200}
				request := fakes.FakeAuthorizedRequest(t, http.MethodPost, uri, accOrigin.CPF.String(), secret.String(), reqBody)

				return request
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "should return 401 if user is not authorized",
			requestSetup: func(t *testing.T) *http.Request {
				accDestination := createAccount(t, testdata.GetValidCPF(), testdata.GetValidSecret(), 0)

				reqBody := schemes.CreateTransferInput{AccountDestinationID: accDestination.ID.String(), Amount: 100}
				request := fakes.FakeRequest(http.MethodPost, uri, reqBody)

				return request
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.DefaultClient.Do(tt.requestSetup(t))
			require.NoError(t, err)

			var body bytes.Buffer

			_, err = io.Copy(&body, resp.Body)
			require.NoError(t, err)

			t.Log(body.String())

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			dockertest.TruncateTables(context.Background(), testenv.DB)
		})
	}
}
