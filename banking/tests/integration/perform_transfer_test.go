package integration

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/banking/domain/vos"
	"github.com/thalissonfelipe/banking/banking/gateway/http/transfer/schema"
	"github.com/thalissonfelipe/banking/banking/tests/dockertest"
	"github.com/thalissonfelipe/banking/banking/tests/fakes"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
	"github.com/thalissonfelipe/banking/banking/tests/testenv"
)

func TestIntegration_PerformTransfer(t *testing.T) {
	uri := testenv.ServerURL + "/api/v1/transfers"

	testCases := []struct {
		name           string
		requestSetup   func(t *testing.T) *http.Request
		expectedStatus int
	}{
		{
			name: "should transfer successfully",
			requestSetup: func(t *testing.T) *http.Request {
				secret := testdata.Secret()
				acc1 := createAccount(t, testdata.CPF(), secret, 100)
				acc2 := createAccount(t, testdata.CPF(), testdata.Secret(), 0)

				reqBody := schema.PerformTransferInput{AccountDestinationID: acc2.ID.String(), Amount: 100}
				request := fakes.FakeAuthorizedRequest(
					t, http.MethodPost, uri, acc1.CPF.String(), secret.String(), reqBody)

				return request
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "should return 400 if account does not has sufficient funds",
			requestSetup: func(t *testing.T) *http.Request {
				secret := testdata.Secret()
				acc1 := createAccount(t, testdata.CPF(), secret, 100)
				acc2 := createAccount(t, testdata.CPF(), testdata.Secret(), 0)

				reqBody := schema.PerformTransferInput{AccountDestinationID: acc2.ID.String(), Amount: 200}
				request := fakes.FakeAuthorizedRequest(
					t, http.MethodPost, uri, acc1.CPF.String(), secret.String(), reqBody)

				return request
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "should return 400 if the account_destination_id is the same as the account origin",
			requestSetup: func(t *testing.T) *http.Request {
				secret := testdata.Secret()
				accOrigin := createAccount(t, testdata.CPF(), secret, 100)

				reqBody := schema.PerformTransferInput{AccountDestinationID: accOrigin.ID.String(), Amount: 200}
				request := fakes.FakeAuthorizedRequest(
					t, http.MethodPost, uri, accOrigin.CPF.String(), secret.String(), reqBody)

				return request
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "should return 404 if account origin does not exist anymore",
			requestSetup: func(t *testing.T) *http.Request {
				secret := testdata.Secret()
				accOrigin := createAccount(t, testdata.CPF(), secret, 100)

				reqBody := schema.PerformTransferInput{AccountDestinationID: vos.NewAccountID().String(), Amount: 200}
				request := fakes.FakeAuthorizedRequest(
					t, http.MethodPost, uri, accOrigin.CPF.String(), secret.String(), reqBody)

				dockertest.TruncateTables(context.Background(), testenv.DB)

				return request
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "should return 404 if account destination does not exist",
			requestSetup: func(t *testing.T) *http.Request {
				secret := testdata.Secret()
				accOrigin := createAccount(t, testdata.CPF(), secret, 100)

				reqBody := schema.PerformTransferInput{AccountDestinationID: vos.NewAccountID().String(), Amount: 200}
				request := fakes.FakeAuthorizedRequest(
					t, http.MethodPost, uri, accOrigin.CPF.String(), secret.String(), reqBody)

				return request
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "should return 401 if user is not authorized",
			requestSetup: func(t *testing.T) *http.Request {
				accDestination := createAccount(t, testdata.CPF(), testdata.Secret(), 0)

				reqBody := schema.PerformTransferInput{AccountDestinationID: accDestination.ID.String(), Amount: 100}
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

			defer resp.Body.Close()

			var body bytes.Buffer

			_, err = io.Copy(&body, resp.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			dockertest.TruncateTables(context.Background(), testenv.DB)
		})
	}
}
