package integration

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/banking/gateways/http/auth/schemes"
	"github.com/thalissonfelipe/banking/banking/tests/dockertest"
	"github.com/thalissonfelipe/banking/banking/tests/fakes"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
	"github.com/thalissonfelipe/banking/banking/tests/testenv"
)

func TestIntegration_Login(t *testing.T) {
	testCases := []struct {
		name           string
		requestSetup   func() *http.Request
		expectedStatus int
	}{
		{
			name: "should log in successfully",
			requestSetup: func() *http.Request {
				secret := testdata.GetValidSecret()
				acc := createAccount(t, testdata.GetValidCPF(), secret, 0)

				reqBody := schemes.LoginInput{CPF: acc.CPF.String(), Secret: secret.String()}
				request := fakes.FakeRequest(http.MethodPost, testenv.ServerURL+"/api/v1/login", reqBody)

				return request
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "should return status code 400 if account does not exist",
			requestSetup: func() *http.Request {
				reqBody := schemes.LoginInput{
					CPF:    testdata.GetValidCPF().String(),
					Secret: testdata.GetValidSecret().String(),
				}
				request := fakes.FakeRequest(http.MethodPost, testenv.ServerURL+"/api/v1/login", reqBody)

				return request
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "should return status code 400 if secret is not valid",
			requestSetup: func() *http.Request {
				acc := createAccount(t, testdata.GetValidCPF(), testdata.GetValidSecret(), 0)

				reqBody := schemes.LoginInput{CPF: acc.CPF.String(), Secret: "12345678"}
				request := fakes.FakeRequest(http.MethodPost, testenv.ServerURL+"/api/v1/login", reqBody)

				return request
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.DefaultClient.Do(tt.requestSetup())
			require.NoError(t, err)

			defer resp.Body.Close()

			var body bytes.Buffer

			_, err = io.Copy(&body, resp.Body)
			require.NoError(t, err)

			t.Log(body.String())

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			dockertest.TruncateTables(context.Background(), testenv.DB)
		})
	}
}
