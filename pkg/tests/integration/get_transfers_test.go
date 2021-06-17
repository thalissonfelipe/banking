package integration

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/pkg/tests/dockertest"
	"github.com/thalissonfelipe/banking/pkg/tests/fakes"
	"github.com/thalissonfelipe/banking/pkg/tests/testdata"
	"github.com/thalissonfelipe/banking/pkg/tests/testenv"
)

func TestIntegration_GetTransfers(t *testing.T) {
	uri := testenv.ServerURL + "/api/v1/transfers"

	testCases := []struct {
		name           string
		requestSetup   func() *http.Request
		expectedStatus int
	}{
		{
			name: "should return list of transfers sucessfully",
			requestSetup: func() *http.Request {
				secret := testdata.GetValidSecret()
				acc := createAccount(t, testdata.GetValidCPF(), secret, 0)

				request := fakes.FakeAuthorizedRequest(t, http.MethodGet, uri, acc.CPF.String(), secret.String(), nil)

				return request
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "should return status code 401 when user is not authorized",
			requestSetup: func() *http.Request {
				request := fakes.FakeRequest(http.MethodGet, testenv.ServerURL+"/api/v1/transfers", nil)

				return request
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.DefaultClient.Do(tt.requestSetup())
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
